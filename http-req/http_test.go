package httpreq_test

import (
	"crypto/tls"
	"testing"

	httpreq "github.com/sdjqwbz/go-mods/http-req"
)

func TestHttpGet(t *testing.T) {
	res, err := httpreq.New("https://www.baidu.com").
		SetTlsVerify(false).
		Get("", nil)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(res))
	}
}

func TestHttpQuickGet(t *testing.T) {
	res, err := httpreq.QuickGet("https://www.baidu.com", nil, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(res))
	}
}
