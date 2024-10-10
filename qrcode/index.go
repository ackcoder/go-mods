package qrcode

import "image"

type QrcodeInfoOption func(qr *QrcodeInfo)

type QrcodeInfo struct {
	bgImg  image.Image //基底二维码图像实例
	color  *image.RGBA //最终带色彩的二维码图像数据
	Detail QrcodeDetail
}

type QrcodeDetail struct {
	content       string
	imgSize       int    //二维码大小
	imgPath       string //二维码保存文件路径、可选，要存文件就必须配置
	centerImg     string //中心图文件路径或Base64字串、可选，要设置中心图就必须配置
	centerImgSize [2]int //中心图宽高、可选，设置了中心图就必须配置
}

// 创建二维码实例
func NewQrcode(options ...QrcodeInfoOption) *QrcodeInfo {
	qr := &QrcodeInfo{
		bgImg:  nil,
		color:  nil,
		Detail: QrcodeDetail{},
	}
	for _, fn := range options {
		fn(qr)
	}
	return qr
}

// 配置二维码内容
func WithContent(txt string) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Detail.content = txt
	}
}

// 配置二维码大小
func WithSize(size int) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Detail.imgSize = size
	}
}

// 配置二维码保存文件路径
func WithSavePath(path string) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Detail.imgPath = path
	}
}

// 配置中心图保存文件路径
func WithCenterImg(pathOrB64str string) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Detail.centerImg = pathOrB64str
	}
}

// 配置中心图大小
func WithCenterSize(w, h int) QrcodeInfoOption {
	return func(qr *QrcodeInfo) {
		qr.Detail.centerImgSize = [2]int{w, h}
	}
}
