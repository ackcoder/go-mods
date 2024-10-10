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

// 设置基础二维码图像
func (qr *QrcodeInfo) SetBaseImage() error {
	var err error
	var qrPkg *qrcodePkg.QRCode
	qrPkg, err = qrcodePkg.New(qr.Detail.content, qrcodePkg.Highest)
	if err != nil {
		return err
	}
	qrPkg.DisableBorder = true
	qr.bgImg = qrPkg.Image(qr.Detail.imgSize) //生成基底二维码

	// 保存二维码数据
	b := qr.bgImg.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, qr.bgImg, image.Point{X: 0, Y: 0}, draw.Src)
	qr.color = m
	return nil
}

// 设置中心图到二维码图像上
func (qr *QrcodeInfo) SetCenterImage() (err error) {
	if qr.Detail.centerImg == "" {
		return errors.New("未配置中心图文件路径")
	}
	if qr.Detail.centerImgSize[0] == 0 || qr.Detail.centerImgSize[1] == 0 {
		return errors.New("未配置中心图大小")
	}
	if qr.Detail.centerImgSize[0] > qr.Detail.imgSize || qr.Detail.centerImgSize[1] > qr.Detail.imgSize {
		return errors.New("中心图大小不能超过基底二维码大小")
	}

	// 尝试获取中心图数据
	var fileHandle io.Reader
	if _, tmpE := os.Stat(qr.Detail.centerImg); tmpE == nil {
		fileHandle, err = os.Open(qr.Detail.centerImg)
		if err != nil {
			return err
		}
	} else {
		b64, err := base64.StdEncoding.DecodeString(qr.Detail.centerImg)
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
		uint(qr.Detail.centerImgSize[0]),
		uint(qr.Detail.centerImgSize[1]),
		imgHandle,
		resize.Lanczos3,
	)

	// 绘制基底二维码及其上的中心图、保存二维码数据
	b := qr.bgImg.Bounds()
	rect := imgHandle.Bounds().Add(image.Pt(
		(b.Max.X-imgHandle.Bounds().Max.X)/2,
		(b.Max.Y-imgHandle.Bounds().Max.Y)/2,
	))
	m := image.NewRGBA(b)
	draw.Draw(m, b, qr.bgImg, image.Point{X: 0, Y: 0}, draw.Src)
	draw.Draw(m, rect, imgHandle, image.Point{X: 0, Y: 0}, draw.Over)
	qr.color = m
	return nil
}

// 二维码存为文件
func (qr *QrcodeInfo) SaveToFile() (err error) {
	if qr.Detail.imgPath == "" {
		return errors.New("未配置二维码保存路径")
	}
	if qr.color == nil {
		return errors.New("未获取到二维码图像数据、请检查 Set* 相关方法是否正确执行")
	}

	file, err := os.Create(qr.Detail.imgPath)
	if err != nil {
		return
	}
	return png.Encode(file, qr.color)
}

// 二维码转为base64字串
//
//	注: 前端按需补充 "data:image/png;base64," 前缀以显示图片
func (qr *QrcodeInfo) SaveAsBase64Str() (str string, err error) {
	if qr.color == nil {
		err = errors.New("未获取到二维码图像数据、请检查 Set* 相关方法是否正确执行")
		return
	}

	buff := bytes.NewBuffer(nil)
	// 图像写入buff
	png.Encode(buff, qr.color)
	// buff转base64字串
	str = base64.StdEncoding.EncodeToString(buff.Bytes())
	return
}
