package httpreq_test

import (
	"testing"

	httpreq "github.com/sdjqwbz/go-mods/http-req"
)

func TestHttpGet(t *testing.T) {
	res, err := httpreq.New("https://www.baidu.com").
		SetTimeout(6).
		SetTlsServerSkipVerify().
		Get("", nil)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(res))
	}
}

func TestHttpQuickGet(t *testing.T) {
	res, err := httpreq.QuickGet("https://www.baidu.com", nil)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(res))
	}
}
