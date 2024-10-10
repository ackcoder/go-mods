package qrcode

import "testing"

func TestQrcode(t *testing.T) {
	qr := NewQrcode(
		WithContent("hello world!"),
		WithSize(120),
		// 中心图配置
		// WithCenterImg(DefaultCenterImage),
		// WithCenterSize(30, 30),
	)
	qr.SetBaseImage()
	// 填入中心图
	// qr.SetCenterImage()
	imgStr, err := qr.SaveAsBase64Str()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s\n", imgStr)
}
