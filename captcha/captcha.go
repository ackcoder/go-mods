package captcha

import (
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	dirver base64Captcha.Driver
	store  base64Captcha.Store

	length int //验证码长度
}

// 创建验证码实例
//   - {size} 验证码长度
//   - {exp} 验证码有效期,单位/秒
//   - {dirverType} 验证码类型,可选,默认(数值+字母组合)
func New(size, exp int, dirverType ...CaptchaType) *Captcha {
	ins := new(Captcha)
	ins.length = size
	ins.store = base64Captcha.NewMemoryStore(
		base64Captcha.GCLimitNumber,
		time.Duration(exp)*time.Second,
	)
	if len(dirverType) != 0 {
		dirverType[0](ins)
	} else {
		WithTypeString(120, 40)(ins)
	}
	return ins
}

// Make 生成验证码
//   - {idKey} 验证码ID 校验时要用
//   - {b64Str} 验证码内容 base64 字串
func (c *Captcha) Make() (idKey, b64Str string, err error) {
	idKey, question, answer := c.dirver.GenerateIdQuestionAnswer()
	if err = c.store.Set(idKey, answer); err != nil {
		return
	}
	item, err := c.dirver.DrawCaptcha(question)
	if err != nil {
		return
	}
	b64Str = item.EncodeB64string()
	return
}

// Check 校验验证码
//   - {idKey} 验证码ID
//   - {code} 用户提交的验证码
func (c *Captcha) Check(idKey, code string) bool {
	return c.store.Verify(idKey, strings.ToUpper(code), true)
}
