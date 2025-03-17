package httpreq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Get 请求
func (rc *ReqClient) Get(api string, headers map[string]string) (res []byte, err error) {
	req, err := baseNewRequest(http.MethodGet, rc.domain+api, headers, "")
	if err != nil {
		return
	}

	if rc.client == nil {
		rc.newClient()
	}
	res, err = baseDoRequest(*rc.client, req)
	return
}

// Post 请求
func (rc *ReqClient) Post(api string, headers, bodys map[string]string) (res []byte, err error) {
	bodyJson, err := json.Marshal(bodys)
	if err != nil {
		return
	}
	req, err := baseNewRequest(http.MethodPost, rc.domain+api, headers, string(bodyJson))
	if err != nil {
		return
	}

	if rc.client == nil {
		rc.newClient()
	}
	res, err = baseDoRequest(*rc.client, req)
	return
}

// Put 请求
func (rc *ReqClient) Put(api string, headers, bodys map[string]string) (res []byte, err error) {
	bodyJson, err := json.Marshal(bodys)
	if err != nil {
		return
	}
	req, err := baseNewRequest(http.MethodPut, rc.domain+api, headers, string(bodyJson))
	if err != nil {
		return
	}

	if rc.client == nil {
		rc.newClient()
	}
	res, err = baseDoRequest(*rc.client, req)
	return
}

// Delete 请求
func (rc *ReqClient) Delete(api string, headers map[string]string) (res []byte, err error) {
	req, err := baseNewRequest(http.MethodDelete, rc.domain+api, headers, "")
	if err != nil {
		return
	}

	if rc.client == nil {
		rc.newClient()
	}
	res, err = baseDoRequest(*rc.client, req)
	return
}

// ======================================================================

// QuickGet 快速请求
func QuickGet(url string, headers map[string]string, tlsConf ...*tls.Config) (res []byte, err error) {
	req, err := baseNewRequest(http.MethodGet, url, headers, "")
	if err != nil {
		return
	}

	var client http.Client
	if len(tlsConf) != 0 {
		trans := http.Transport{TLSClientConfig: tlsConf[0]}
		client = http.Client{Transport: &trans}
	} else {
		client = http.Client{}
	}

	res, err = baseDoRequest(client, req)
	return
}

// QuickPost 快速请求
func QuickPost(url string, headers, bodys map[string]string, tlsConf ...*tls.Config) (res []byte, err error) {
	bodyJson, err := json.Marshal(bodys)
	if err != nil {
		return
	}
	req, err := baseNewRequest(http.MethodPost, url, headers, string(bodyJson))
	if err != nil {
		return
	}

	var client http.Client
	if len(tlsConf) != 0 {
		trans := http.Transport{TLSClientConfig: tlsConf[0]}
		client = http.Client{Transport: &trans}
	} else {
		client = http.Client{}
	}

	res, err = baseDoRequest(client, req)
	return
}

// QuickPut 快速请求
func QuickPut(url string, headers, bodys map[string]string, tlsConf ...*tls.Config) (res []byte, err error) {
	bodyJson, err := json.Marshal(bodys)
	if err != nil {
		return
	}
	req, err := baseNewRequest(http.MethodPut, url, headers, string(bodyJson))
	if err != nil {
		return
	}

	var client http.Client
	if len(tlsConf) != 0 {
		trans := http.Transport{TLSClientConfig: tlsConf[0]}
		client = http.Client{Transport: &trans}
	} else {
		client = http.Client{}
	}

	res, err = baseDoRequest(client, req)
	return
}

// QuickDelete 快速请求
func QuickDelete(url string, headers map[string]string, tlsConf ...*tls.Config) (res []byte, err error) {
	req, err := baseNewRequest(http.MethodDelete, url, headers, "")
	if err != nil {
		return
	}

	var client http.Client
	if len(tlsConf) != 0 {
		trans := http.Transport{TLSClientConfig: tlsConf[0]}
		client = http.Client{Transport: &trans}
	} else {
		client = http.Client{}
	}

	res, err = baseDoRequest(client, req)
	return
}

// ======================================================================

// 基础 request 封装
func baseNewRequest(method, url string, headers map[string]string, body string) (req *http.Request, err error) {
	switch method {
	case http.MethodGet:
		fallthrough
	case http.MethodDelete:
		req, err = http.NewRequest(method, url, nil)
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		var read *bytes.Reader
		if body != "" {
			read = bytes.NewReader([]byte(body))
		}
		req, err = http.NewRequest(method, url, read)
	default:
		err = errors.New("暂不支持的请求方法")
	}
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return
}

// 基础 request 执行封装
func baseDoRequest(c http.Client, req *http.Request) (res []byte, err error) {
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	res, err = io.ReadAll(resp.Body)
	return
}
