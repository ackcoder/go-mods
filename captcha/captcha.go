package captcha

import (
	"image/color"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/mojocn/base64Captcha"
)

// TODO 目前简单引入、后续向qrcode做类似改造
// TODO 两个结构体方法都new实例、需要改进

type Captcha struct {
	dirver *base64Captcha.DriverString
	store  base64Captcha.Store
}

func New() *Captcha {
	return &Captcha{
		store: base64Captcha.DefaultMemStore, //内存存储,默认验证码有效期10分钟、校验后删除
		dirver: &base64Captcha.DriverString{
			Height:          60,
			Width:           200,
			NoiseCount:      32,    //噪点数
			ShowLineOptions: 2 | 4, //干扰线
			Length:          4,
			Source:          "123456789ABCDEFGHIJKLHIJKLMNOPQRSTUVWXYZ",
			BgColor:         &color.RGBA{R: 10, G: 20, B: 50, A: 10},
			Fonts:           []string{"wqy-microhei.ttc"},
		},
	}
}

// Make 生成验证码
//   - {idKey} 验证码ID 校验时要用
//   - {b64img} 验证码图片 base64 字串
func (c *Captcha) Make() (idKey, b64img string, err error) {
	d := c.dirver.ConvertFonts() //加载字体
	ins := base64Captcha.NewCaptcha(d, c.store)
	idKey, b64img, _, err = ins.Generate()
	return
}

// Check 校验验证码
//   - {idKey} 验证码ID
//   - {code} 用户提交的验证码
func (c *Captcha) Check(idKey, code string) bool {
	ins := base64Captcha.NewCaptcha(c.dirver, c.store)
	code = gstr.ToUpper(code)
	return ins.Verify(idKey, code, true) //校验后删除验证码缓存
}
