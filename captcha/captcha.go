package captcha

import (
	"image/color"
	"strings"

	"github.com/mojocn/base64Captcha"
)

type CaptchaOption func(c *Captcha)

type Captcha struct {
	dirver *base64Captcha.DriverString
}

func New(options ...CaptchaOption) *Captcha {
	ins := new(Captcha)
	ins.dirver = new(base64Captcha.DriverString)
	// 默认值
	WithCodeNum(4)(ins)
	WithSize(100, 35)(ins)
	WithNoise(32)(ins)
	WithLine(2)(ins)
	WithBgColor(10, 20, 50, 10)(ins)
	WithSource("123456789ABCDEFGHIJKLHIJKLMNOPQRSTUVWXYZ")(ins)
	WithFont("wqy-microhei.ttc")(ins)
	// 用户配置
	for _, fn := range options {
		fn(ins)
	}
	return ins
}

// 配置验证码字符数
//   - {num} 字符数,默认4
func WithCodeNum(num uint) CaptchaOption {
	return func(c *Captcha) {
		c.dirver.Length = int(num)
	}
}

// 配置验证码大小
//   - {w},{h} 宽高,单位/像素px  默认100x35
func WithSize(w, h uint) CaptchaOption {
	return func(c *Captcha) {
		c.dirver.Width = int(w)
		c.dirver.Height = int(h)
	}
}

// 配置验证码噪点数
//   - {num} 噪点数,默认32
func WithNoise(num uint) CaptchaOption {
	return func(c *Captcha) {
		c.dirver.NoiseCount = int(num)
	}
}

// 配置验证码干扰线
//   - {num} 干扰线,默认2,范围0-3,值越大干扰越重
func WithLine(num uint) CaptchaOption {
	return func(c *Captcha) {
		switch num {
		case 0:
			c.dirver.ShowLineOptions = 0
		case 1:
			c.dirver.ShowLineOptions = 2
		case 3:
			c.dirver.ShowLineOptions = 2 | 4 | 8
		case 2:
			fallthrough
		default:
			c.dirver.ShowLineOptions = 2 | 4
		}
	}
}

// 配置验证码背景色
//   - {r},{g},{b},{a} 颜色值,默认10,20,50,10
func WithBgColor(r, g, b uint8, a ...uint8) CaptchaOption {
	return func(c *Captcha) {
		var realA uint8 = 255
		if len(a) > 0 {
			realA = a[0]
		}
		c.dirver.BgColor = &color.RGBA{R: r, G: g, B: b, A: realA}
	}
}

// 配置验证码取值源
func WithSource(source string) CaptchaOption {
	return func(c *Captcha) {
		c.dirver.Source = source
	}
}

// 配置验证码字体
func WithFont(font string) CaptchaOption {
	return func(c *Captcha) {
		c.dirver.Fonts = []string{font}
	}
}

// Make 生成验证码
//   - {idKey} 验证码ID 校验时要用
//   - {b64img} 验证码图片 base64 字串
func (c *Captcha) Make() (idKey, b64img string, err error) {
	idKey, b64img, _, err = base64Captcha.NewCaptcha(
		c.dirver.ConvertFonts(),       //加载字体
		base64Captcha.DefaultMemStore, //内存存储,默认有效期10分钟、校验后删除
	).Generate()
	return
}

// Check 校验验证码
//   - {idKey} 验证码ID
//   - {code} 用户提交的验证码
func (c *Captcha) Check(idKey, code string) bool {
	return base64Captcha.
		NewCaptcha(c.dirver, base64Captcha.DefaultMemStore).
		Verify(idKey, strings.ToUpper(code), true) //校验后删除缓存
}
