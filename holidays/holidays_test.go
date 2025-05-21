package holidays_test

import (
	"testing"
	"time"

	"github.com/ackcoder/go-mods/holidays"
)

func TestGet(t *testing.T) {
	_, err := holidays.Get(-1)
	if err != nil {
		t.Error(err)
	}

	res, err := holidays.Get(2024)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", res)

	time.Sleep(time.Second)

	// repeat
	res, err = holidays.Get(2024)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", res)
}
