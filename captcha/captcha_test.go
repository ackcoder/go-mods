package captcha_test

import (
	"testing"

	"github.com/sdjqwbz/go-mods/captcha"
)

func TestCaptcha(t *testing.T) {
	ins := captcha.New()
	tk, b64Img, err := ins.Make()
	if err != nil {
		t.Error(err)
	}
	t.Log(tk, b64Img)

	if ins.Check(tk, "wrong_code") {
		t.Log("验证通过")
	} else {
		t.Error("验证失败")
	}
}
