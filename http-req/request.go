package httpreq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sdjqwbz/go-mods/utils"
)

// Get 请求
func (rc *ReqClient) Get(api string, headers map[string]string) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	req := baseNewRequest(http.MethodGet, rc.domain+api, headers, "")
	res = baseDoRequest(rc.client, req)
	return
}

// Post 请求
func (rc *ReqClient) Post(api string, headers, bodys map[string]string) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	bodyJson, err := json.Marshal(bodys)
	utils.Must(err, "格式化请求参数异常")
	req := baseNewRequest(http.MethodPost, rc.domain+api, headers, string(bodyJson))
	res = baseDoRequest(rc.client, req)
	return
}

// Put 请求
func (rc *ReqClient) Put(api string, headers, bodys map[string]string) (res []byte, err error) {
	return rc.Post(api, headers, bodys)
}

// QuickGet 快速请求
func QuickGet(url string, headers map[string]string, tlsConf ...*tls.Config) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	req := baseNewRequest(http.MethodGet, url, headers, "")

	var client http.Client
	if len(tlsConf) != 0 {
		client = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConf[0],
			},
		}
	} else {
		client = http.Client{}
	}

	res = baseDoRequest(client, req)
	return
}

// QuickPost 快速请求
func QuickPost(url string, headers, bodys map[string]string, tlsConf ...*tls.Config) (res []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	bodyJson, err := json.Marshal(bodys)
	utils.Must(err, "格式化请求参数异常")
	req := baseNewRequest(http.MethodPost, url, headers, string(bodyJson))

	var client http.Client
	if len(tlsConf) != 0 {
		client = http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConf[0],
			},
		}
	} else {
		client = http.Client{}
	}

	res = baseDoRequest(client, req)
	return
}

// QuickPut 快速请求
func QuickPut(url string, headers, bodys map[string]string, tlsConf ...*tls.Config) (res []byte, err error) {
	return QuickPost(url, headers, bodys, tlsConf...)
}

// 基础 request 封装 (含panic)
func baseNewRequest(method, url string, headers map[string]string, body string) (req *http.Request) {
	var err error
	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(http.MethodGet, url, nil)
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		var read *bytes.Reader
		if body != "" {
			read = bytes.NewReader([]byte(body))
		}
		req, err = http.NewRequest(method, url, read)
	default:
		panic("暂不支持的请求方法")
	}
	utils.Must(err, "创建Http请求异常")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return
}

// 基础 request 执行封装 (含panic)
func baseDoRequest(c http.Client, req *http.Request) (res []byte) {
	resp, err := c.Do(req)
	utils.Must(err, "执行Http请求异常")
	defer resp.Body.Close()

	res, err = io.ReadAll(resp.Body)
	utils.Must(err, "读取Http响应异常")
	return
}
