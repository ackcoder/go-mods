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

// SaveToFile 二维码存为文件
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

// SaveAsBase64Str 二维码转为base64字串
func (qr *Qrcode) SaveAsBase64Str() (str string, err error) {
	if err = qr.drawImage(); err != nil {
		return
	}

	buff := bytes.NewBuffer(nil)
	if err = png.Encode(buff, qr.color); err != nil {
		return
	}
	str = base64.StdEncoding.EncodeToString(buff.Bytes())
	return
}

// SaveAsBase64Img 二维码转为base64图片字串
func (qr *Qrcode) SaveAsBase64Img() (str string, err error) {
	str, err = qr.SaveAsBase64Str()
	str = "data:image/png;base64," + str
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
