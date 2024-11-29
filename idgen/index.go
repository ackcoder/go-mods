package idgen

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
