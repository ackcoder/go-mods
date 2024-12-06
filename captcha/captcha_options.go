package captcha

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
)

// 验证码类型项
type CaptchaType func(c *Captcha)

// 图片类验证码选项参数
type ImageOption struct {
	Noise     int         //噪点数,默认32
	LineLevel int         //干扰线强度,范围1-3,-1表示不设置,默认2
	BgColor   *color.RGBA //背景色,默认r10,g20,b50,a10
	Source    string      //取值源
	Font      string      //字体
}

// 音频验证码
//   - {lang} 音频语言,可用值:"en"英语,"ja"日语,"ru"俄语,"zh"中文,默认"en"
func SetAudio(lang ...string) CaptchaType {
	return func(c *Captcha) {
		var language string
		if len(lang) != 0 {
			language = lang[0]
		} else {
			language = "en"
		}
		c.dirver = base64Captcha.NewDriverAudio(c.length, language)
	}
}

// 中文验证码
//   - {w},{h} 宽高
//   - {opt} 可选项参数
func SetChinese(w, h int, opt ...*ImageOption) CaptchaType {
	return func(c *Captcha) {
		noise, line, bg, source, font := takeImageOptionValues(opt...)
		if source == "" {
			source = base64Captcha.TxtChineseCharaters
		}
		if font == "" {
			font = "wqy-microhei.ttc"
		}
		c.dirver = base64Captcha.NewDriverChinese(
			h, w,
			noise, line, c.length, source, bg,
			base64Captcha.DefaultEmbeddedFonts,
			[]string{font},
		)
	}
}

// 纯数字验证码
//   - {w},{h} 宽高,最低100x36
func SetDigit(w, h int) CaptchaType {
	return func(c *Captcha) {
		if w < 100 {
			w = 100
		}
		if h < 36 {
			h = 36
		}
		c.dirver = base64Captcha.NewDriverDigit(h, w, c.length, 0.8, 26)
	}
}

// 数学计算验证码
//   - {w},{h} 宽高
//   - {opt} 可选项参数
func SetMath(w, h int, opt ...*ImageOption) CaptchaType {
	return func(c *Captcha) {
		noise, line, bg, _, font := takeImageOptionValues(opt...)
		if font == "" {
			font = "wqy-microhei.ttc"
		}
		// 注: 不使用字体
		// c.dirver = base64Captcha.NewDriverMath(
		// 	h, w, noise, line, bg,
		// 	base64Captcha.DefaultEmbeddedFonts,
		// 	[]string{font},
		// ).ConvertFonts()
		c.dirver = base64Captcha.NewDriverMath(
			h, w, noise, line, bg, nil, nil,
		)
	}
}

// 数值字母验证码
//   - {w},{h} 宽高
//   - {opt} 可选项参数
func SetString(w, h int, opt ...*ImageOption) CaptchaType {
	return func(c *Captcha) {
		noise, line, bg, source, font := takeImageOptionValues(opt...)
		if source == "" {
			source = base64Captcha.TxtSimpleCharaters
		}
		if font == "" {
			font = "wqy-microhei.ttc"
		}
		c.dirver = base64Captcha.NewDriverString(
			h, w,
			noise, line, c.length, source, bg,
			base64Captcha.DefaultEmbeddedFonts,
			[]string{font},
		).ConvertFonts()
	}
}

func takeImageOptionValues(opt ...*ImageOption) (
	noise, line int,
	bg *color.RGBA,
	source, font string,
) {
	var currOpt = ImageOption{}
	if len(opt) != 0 {
		currOpt = *opt[0]
	}
	// 默认值配置
	if currOpt.Noise == 0 {
		noise = 32
	}
	if currOpt.LineLevel == 0 {
		currOpt.LineLevel = 2
	}
	if currOpt.BgColor == nil {
		currOpt.BgColor = &color.RGBA{R: 10, G: 20, B: 50, A: 10}
	}
	// 设置返回值
	noise = currOpt.Noise
	bg = currOpt.BgColor
	source = currOpt.Source
	font = currOpt.Font
	switch currOpt.LineLevel {
	case -1:
		line = 0
	case 1:
		line = 2
	case 3:
		line = 2 | 4 | 8
	case 2:
		fallthrough
	default:
		line = 2 | 4
	}
	return
}
