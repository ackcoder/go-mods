package captcha_test

import (
	"testing"

	"github.com/sdjqwbz/go-mods/captcha"
)

func TestCaptcha(t *testing.T) {
	// 1. 默认为数值字母验证码
	// ins := captcha.New(4, 60) //默认宽高120x40
	// ins := captcha.New(4, 60, captcha.WithTypeString(150, 46))
	// 2. 语音验证码
	// ins := captcha.New(4, 60, captcha.WithTypeAudio("zh"))
	// 3. 中文验证码
	ins := captcha.New(5, 72, captcha.WithTypeChinese(100, 35, &captcha.ImageOption{
		Noise:     0,
		LineLevel: -1,
	}))
	// 4. 数字验证码
	// ins := captcha.New(4, 60, captcha.WithTypeDigit(100, 37))
	// 5. 数学计算验证码
	// ins := captcha.New(4, 60, captcha.WithTypeMath(120, 40))

	tk, b64Str, err := ins.Make()
	if err != nil {
		t.Error(err)
	}
	t.Log(tk, b64Str)

	if ins.Check(tk, "wrong_code") {
		t.Log("验证通过")
	} else {
		t.Error("验证失败")
	}
}
