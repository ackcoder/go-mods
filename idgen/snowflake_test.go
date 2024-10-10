package idgen

import (
	"testing"
)

func TestIdGenerator_Gen(t *testing.T) {
	// test in goroutine
	for i := 0; i < 100; i++ {
		go Init(1, 1)
	}

	for i := 0; i < 100; i++ {
		go func() {
			id, _ := GenB36()
			t.Logf("%s\n", id)
		}()
	}
}
