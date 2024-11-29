package idgen_test

import (
	"testing"
	"time"

	"github.com/sdjqwbz/go-mods/idgen"
)

func TestIdGenerator_New(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()
	// idgen.New(1, 32) //wrong
	idgen.New(1, 1, time.Now().Add(1*time.Hour)) //no wrong
}

func TestIdGenerator(t *testing.T) {
	id := idgen.New(1, 1)
	for i := 0; i < 100; i++ {
		go func() {
			t.Log(id.Gen())
		}()
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
