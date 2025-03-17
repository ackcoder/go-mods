package idgen

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// idGenerator 雪花算法 id 生成器
// 原理: https://en.wikipedia.org/wiki/Snowflake_ID
// 原作者: 小生凡一
// 参考文章:
// 美团Leaf方案 https://tech.meituan.com/2017/04/21/mt-leaf.html
// 百度UidGenerator方案 https://zhuanlan.zhihu.com/p/550596015

const (
	// timestamp occupancy bits 时间戳占用位
	timestampBits = 41
	// dataCenterId occupancy bits 集群序号占用位
	dataCenterIdBits = 5
	// workerId occupancy bits 程序序号占用位
	workerIdBits = 5
	// sequence occupancy bits 序列占用位
	sequenceBits = 12

	// timestamp 最大值 (2^41-1 = 2199023255551)
	timestampMaxValue = -1 ^ (-1 << timestampBits)
	// dataCenterId 最大值 (2^5-1 = 31)
	dataCenterIdMaxValue = -1 ^ (-1 << dataCenterIdBits)
	// workId 最大值 (2^5-1 = 31)
	workerIdMaxValue = -1 ^ (-1 << workerIdBits)
	// sequence 最大值 (2^12-1 = 4095)
	sequenceMaxValue = -1 ^ (-1 << sequenceBits)

	// workId 偏移位数 (sequenceBits)
	workIdShift = 12
	// dataCenterId 偏移位数 (sequenceBits + workerIdBits)
	dataCenterIdShift = 17
	// timestamp 偏移位数 (sequenceBits + workerIdBits + dataCenterIdBits)
	timestampShift = 22

	defaultInitValue = 0
)

type IdGenerator struct {
	epoch        int64  //起始点 毫秒时间戳, 默认 UTC: 2024-01-01 00:00:00
	timestamp    int64  //记录点 毫秒时间戳 (取 2^41 约 69 年)
	dataCenterId uint64 //机器码 集群ID (取 2^5)
	workerId     uint64 //机器码 程序ID (取 2^5, 机器共 2^10=1024 台)
	sequence     uint64 //序列号 (取 2^12, 即 4096 个)

	mu *sync.Mutex
}

// 创建ID生成器实例
//   - {dataCenterId} 集群ID [0, 31]
//   - {workerId} 程序ID [0, 31]
//   - {start} 可选, 设置起始点, 未来时间会置为默认 (默认为 UTC: 2024-01-01 00:00:00)
func New(dataCenterId, workerId uint64, start ...time.Time) *IdGenerator {
	if dataCenterId > dataCenterIdMaxValue {
		panic(fmt.Sprintf("雪花算法 id 生成器 dataCenterId 范围应为 [0, %d]", dataCenterIdMaxValue))
	}
	if workerId > workerIdMaxValue {
		panic(fmt.Sprintf("雪花算法 id 生成器 workId 范围应为 [0, %d]", workerIdMaxValue))
	}
	realEpoch := time.Date(2024, time.January, 01, 00, 00, 00, 00, time.UTC).UnixMilli()
	if len(start) != 0 && start[0].Before(time.Now()) {
		realEpoch = start[0].UnixMilli()
	}
	return &IdGenerator{
		epoch:        realEpoch,
		timestamp:    defaultInitValue - 1,
		sequence:     defaultInitValue,
		dataCenterId: dataCenterId,
		workerId:     workerId,

		mu: new(sync.Mutex),
	}
}

// Gen 生成ID (10进制位)
func (ig *IdGenerator) Gen() string {
	return fmt.Sprintf("%d", ig.genId())
}

// GenB16 生成雪花算法ID (16进制位, 10进制位)
func (ig *IdGenerator) GenB16() (string, string) {
	id := ig.genId()
	return fmt.Sprintf("%x", id), fmt.Sprintf("%d", id)
}

// GenB36 生成雪花算法ID (36进制位, 10进制位)
//
//	注: 标准库最多支持36位
func (ig *IdGenerator) GenB36() (string, string) {
	id := ig.genId()
	return strconv.FormatUint(id, 36), fmt.Sprintf("%d", id)
}

// GenNum 生成雪花算法ID (10进制位数值)
//
//	注: snowflake设计用到63位二进制位、需要数值必须uint64类型表示
//	    但需注意存数据库/传前端能否正确表示、业务逻辑中也尽量不转非uint64类型
func (ig *IdGenerator) GenNum() uint64 {
	return ig.genId()
}

func (ig *IdGenerator) genId() uint64 {
	ig.mu.Lock()
	defer ig.mu.Unlock()

	var now = time.Now().UnixMilli()

	// 时钟回拨处理
	if ig.timestamp > now {
		// 方式一: 直接抛异常
		// panic(fmt.Sprintf("雪花算法 时钟回拨 最后时间戳 %d, 比较时间戳 %d", ig.timestamp, now))

		// 方式二: 延迟等待三次
		for i := 0; i < 3; i++ {
			time.Sleep(time.Millisecond * 300) //期望时钟自身校正
			now = time.Now().UnixMilli()
			if ig.timestamp <= now {
				break
			}
		}
		if ig.timestamp > now {
			panic(fmt.Sprintf("雪花算法 时钟回拨 最后时间戳 %d, 比较时间戳 %d", ig.timestamp, now))
		}
	}

	if ig.timestamp == now {
		// 相同时间戳、序列号自旋
		ig.sequence = (ig.sequence + 1) & sequenceMaxValue //递增序列号
		if ig.sequence == 0 {
			// 序列号溢出4095、等待至下一毫秒
			for now <= ig.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 时间戳进位、序列号重置
		ig.sequence = defaultInitValue
	}

	diff := uint64(now - ig.epoch)
	if diff > timestampMaxValue {
		// 运行超 69 年期限、直接抛异常
		panic(fmt.Sprintf("雪花算法 起始时间 epoch 范围应为 [0, %d]", timestampMaxValue-1))
	}
	ig.timestamp = now

	return (diff << timestampShift) |
		(ig.dataCenterId << dataCenterIdShift) |
		(ig.workerId << workIdShift) |
		ig.sequence
}
