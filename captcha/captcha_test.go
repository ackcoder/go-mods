package captcha_test

import (
	"testing"

	"github.com/sdjqwbz/go-mods/captcha"
)

func TestCaptcha(t *testing.T) {
	ins := captcha.New(
		4, 60,
		// 1. 测试语音验证码
		// captcha.SetAudio("ja"),
		// 2. 测试中文验证码
		// captcha.SetChinese(100, 35, &captcha.ImageOption{
		// 	Noise:     0,
		// 	LineLevel: -1,
		// }),
		// 3. 测试数字验证码
		// captcha.SetDigit(100, 37),
		// 4. 测试数学计算验证码
		// captcha.SetMath(120, 40),
		// 5. 默认为 数值字母验证码(120x40)
	)
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
