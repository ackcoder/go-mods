package utils_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ackcoder/go-mods/utils"
)

func TestHashPwd_Bcrypt(t *testing.T) {
	const pwd = "a0823h4apowe"

	t.Run("DefaultParams", func(t *testing.T) {
		hash, err := utils.Pwd.MakeByBcrypt(pwd)
		if err != nil {
			t.Fatal("生成密码哈希失败", err)
		}
		t.Log("当前密码哈希", hash)
		if !utils.Pwd.CheckByBcrypt(pwd, hash) {
			t.Fatal("正确密码理应通过验证")
		}
	})

	t.Run("TimeUsage", func(t *testing.T) {
		st := time.Now()
		hash, err := utils.Pwd.MakeByBcrypt(pwd)
		t.Log("生成用时", time.Since(st))
		t.Log("生成结果", hash, err)
	})
}

func TestHashPwd_Scrypt(t *testing.T) {
	const pwd = "ao2ih4nrasd"

	t.Run("DefaultParams", func(t *testing.T) {
		hash, err := utils.Pwd.MakeByScrypt(pwd)
		if err != nil {
			t.Fatal("生成密码哈希失败", err)
		}
		if hash == "" {
			t.Fatal("生成密码哈希为空")
		}
		t.Log("当前密码哈希", hash)
		if !utils.Pwd.CheckByScrypt(pwd, hash) {
			t.Fatal("正确密码理应通过验证")
		}
		if utils.Pwd.CheckByScrypt("xxxxx", hash) {
			t.Fatal("错误密码理应无法通过验证")
		}
	})

	t.Run("TimeUsage", func(t *testing.T) {
		// TODO: 似乎修改 N 值并不影响密码哈希生成耗时？
		utils.Pwd.EditScryptParam(2048, 8, 1, 32, 32)
		st := time.Now()
		hash, err := utils.Pwd.MakeByScrypt(pwd)
		t.Log("生成用时", time.Since(st))
		t.Log("生成结果", hash, err)
	})
}

func TestHashPwd_Argon2(t *testing.T) {
	const pwd = "9uw3j4lkifosi"

	t.Run("DefaultParams", func(t *testing.T) {
		hash, err := utils.Pwd.MakeByArgon2(pwd)
		if err != nil {
			t.Fatal("生成密码哈希失败", err)
		}
		if hash == "" {
			t.Fatal("生成密码哈希为空")
		}
		t.Log("当前密码哈希", hash)
		if !utils.Pwd.CheckByArgon2(pwd, hash) {
			t.Fatal("正确密码理应通过验证")
		}
		if utils.Pwd.CheckByArgon2("xxxxx", hash) {
			t.Fatal("错误密码理应无法通过验证")
		}
	})

	t.Run("TimeUsage", func(t *testing.T) {
		utils.Pwd.EditArgon2Param(1, 64*1024, 4, 32, 32)
		st := time.Now()
		hash, err := utils.Pwd.MakeByArgon2(pwd)
		t.Log("生成用时", time.Since(st))
		t.Log("生成结果", hash, err)
	})
}

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
