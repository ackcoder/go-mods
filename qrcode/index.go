package qrcode

import "image"

type QrcodeOption func(qr *Qrcode)

type Qrcode struct {
	imgContent    string  //二维码内容
	imgSize       uint    //二维码大小,单位/像素px
	centerImg     string  //中心图文件路径或Base64字串
	centerImgSize [2]uint //中心图宽高,单位/像素px

	bgImg image.Image //基底二维码图像实例
	color *image.RGBA //最终带色彩的二维码图像数据
}

// 创建二维码实例
//   - {content} 二维码内容
//   - {options} 可选项, 调用 With* 函数可设置
func New(content string, options ...QrcodeOption) *Qrcode {
	qr := new(Qrcode)
	qr.imgContent = content
	qr.imgSize = 50
	for _, fn := range options {
		fn(qr)
	}
	return qr
}

// 配置二维码大小
//   - {size} 尺寸,单位/像素. 默认50px
func WithSize(size uint) QrcodeOption {
	return func(qr *Qrcode) {
		if size != 0 {
			qr.imgSize = size
		}
	}
}

// 配置中心图
//   - {pathOrB64str} 中心图路径或Base64字串
//   - {w},{h} 中心图宽高,单位/像素
func WithCenterImg(pathOrB64str string, w, h uint) QrcodeOption {
	return func(qr *Qrcode) {
		qr.centerImg = pathOrB64str
		if w != 0 && h != 0 {
			qr.centerImgSize = [2]uint{w, h}
		}
	}
}
