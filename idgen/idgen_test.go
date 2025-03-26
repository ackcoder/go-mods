package idgen_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ackcoder/go-mods/idgen"
)

func TestIdGenerator_New(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()
	// idgen.New(1, 32) //报错
	idgen.New(1, 1, time.Now().Add(1*time.Hour)) //正常执行
}

func TestIdGenerator(t *testing.T) {
	var wg sync.WaitGroup
	start := make(chan struct{}) //同步启动信号,受限于系统调度与CPU数,不一定全部协程都能同时启动

	var checkS []uint64

	id := idgen.New(1, 1)
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			tmpId := id.GenNum()
			checkS = append(checkS, tmpId)
			t.Log(tmpId)
		}()
	}
	time.Sleep(200 * time.Millisecond) //确保所有协程都设置好

	close(start)
	wg.Wait()

	var checkMap = make(map[uint64]struct{})
	for _, v := range checkS {
		if _, ok := checkMap[v]; ok {
			t.Error("id 重复: ", v)
		}
		checkMap[v] = struct{}{}
	}
}

// go test -run='^$' -bench=. -count=1 -benchtime=2s
func BenchmarkIdGenerator(b *testing.B) {
	id := idgen.New(1, 2)
	for i := 0; i < b.N; i++ {
		// 注: 基准测试显示 每次耗时 245ns 左右
		//     正好接近每毫秒每机器 sequence=4096 上限
		id.Gen()
	}
}
