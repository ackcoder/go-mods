package captcha_test

import (
	"testing"

	"github.com/ackcoder/go-mods/captcha"
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

	tk, b64Str, err := ins.Make(false)
	if err != nil {
		t.Error(err)
	}
	t.Log(tk)
	t.Log(b64Str)

	if ins.Check(tk, "wrong_code") {
		t.Error("逻辑错误、应校验失败")
	} else {
		t.Log("逻辑正确、错误输入校验不通过")
	}
}
