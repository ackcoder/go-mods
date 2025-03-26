package qrcode

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"

	"github.com/nfnt/resize"
	qrcodePkg "github.com/skip2/go-qrcode"
)

type QrcodeOption func(qr *Qrcode)

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

// SaveToFile 二维码存为图片文件
//   - {imgPath} 保存文件路径
func (qr *Qrcode) SaveToFile(imgPath string) (err error) {
	if imgPath == "" {
		return errors.New("未配置二维码保存路径")
	}
	if err = qr.drawImage(); err != nil {
		return
	}

	file, err := os.Create(imgPath)
	if err != nil {
		return
	}
	return png.Encode(file, qr.color)
}

// SaveAsB64Str 二维码转为图片 base64 字串
//   - {withFormatPrefix} 是否携带图片格式前缀
func (qr *Qrcode) SaveAsB64Str(withFormatPrefix bool) (str string, err error) {
	if err = qr.drawImage(); err != nil {
		return
	}

	buff := bytes.NewBuffer(nil)
	if err = png.Encode(buff, qr.color); err != nil {
		return
	}
	str = base64.StdEncoding.EncodeToString(buff.Bytes())
	if withFormatPrefix {
		str = "data:image/png;base64," + str
	}
	return
}

// 绘制二维码图像
func (qr *Qrcode) drawImage() error {
	var err error
	var qrPkg *qrcodePkg.QRCode
	qrPkg, err = qrcodePkg.New(qr.imgContent, qrcodePkg.Highest)
	if err != nil {
		return err
	}
	qrPkg.DisableBorder = true              // 去除二维码边框
	qr.bgImg = qrPkg.Image(int(qr.imgSize)) //生成基底二维码

	// 保存二维码数据
	b := qr.bgImg.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, qr.bgImg, image.Point{X: 0, Y: 0}, draw.Src)
	qr.color = m

	// 判断是否需要绘制中心图
	if qr.centerImg == "" {
		return nil
	}
	if qr.centerImgSize[0] > qr.imgSize || qr.centerImgSize[1] > qr.imgSize {
		return errors.New("中心图大小不能超过基底二维码大小")
	}

	// 尝试获取中心图数据
	var fileHandle io.Reader
	if _, tmpE := os.Stat(qr.centerImg); tmpE == nil {
		fileHandle, err = os.Open(qr.centerImg)
		if err != nil {
			return err
		}
	} else {
		b64, err := base64.StdEncoding.DecodeString(qr.centerImg)
		if err != nil {
			return errors.New("获取中心图数据失败、配置项 centerImg 非文件路径或Base64字串")
		}
		fileHandle = bytes.NewBuffer(b64)
	}

	// 解码图像并调整大小
	imgHandle, err := png.Decode(fileHandle)
	if err != nil {
		return err
	}
	imgHandle = resize.Resize(
		qr.centerImgSize[0],
		qr.centerImgSize[1],
		imgHandle,
		resize.Lanczos3,
	)

	// 绘制基底二维码及其上的中心图、保存二维码数据
	rect := imgHandle.Bounds().Add(image.Pt(
		(b.Max.X-imgHandle.Bounds().Max.X)/2,
		(b.Max.Y-imgHandle.Bounds().Max.Y)/2,
	))
	draw.Draw(m, rect, imgHandle, image.Point{X: 0, Y: 0}, draw.Over)
	qr.color = m
	return nil
}
