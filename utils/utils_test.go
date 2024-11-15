package utils_test

import (
	"errors"
	"testing"

	"github.com/sdjqwbz/go-mods/utils"
)

func TestGoVer(t *testing.T) {
	t.Log(utils.GoVer())
	t.Log(utils.GoVerEq(1, 19))
}

func TestNetworkInfoList(t *testing.T) {
	infos, err := utils.NetworkInfoList()
	if err != nil {
		t.Error(err)
	} else {
		for _, info := range infos {
			t.Log(info)
		}
	}
}

func TestGenerateStr(t *testing.T) {
	t.Log(utils.RandStr(10))
	t.Log(utils.RandNumStr(8))
	t.Log(utils.Md5Str16("hello world!"))
}

func TestMustxxx(t *testing.T) {
	defer func() {
		if r := recover(); r != nil && r == "custom error" {
			t.Log("success panic")
		} else {
			t.Error("fail panic")
		}
	}()

	var errEmpty error
	var err = errors.New("custom error")
	var numZero uint32
	var num = 12

	utils.MustZeroN(num)    //pass
	utils.MustZero(numZero) //pass
	utils.MustTrue(true)    //pass
	utils.Must(errEmpty)    //pass
	utils.Must(err)         //panic
}
