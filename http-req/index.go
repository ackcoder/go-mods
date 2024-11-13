package httpreq

import (
	"crypto/tls"
	"net/http"
)

type ReqClient struct {
	tlsConf *tls.Config
	client  http.Client

	domain string //请求目的域, 如"http://xxx.com"
}

func New(domain string) *ReqClient {
	tlsConf := &tls.Config{}
	return &ReqClient{
		domain:  domain,
		tlsConf: tlsConf,
		client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConf,
			},
		},
	}
}

// SetTLSVerify 设置是否校验证书
//
// 注: https自签名证书要么不校验、要么提供客户端证书
func (rc *ReqClient) SetTlsVerify(isCheck bool) *ReqClient {
	rc.tlsConf.InsecureSkipVerify = isCheck
	return rc
}

// SetTlsCertFiles 设置客户端证书 (可能失败并返回nil)
//
// 注: https自签名证书要么不校验、要么提供客户端证书
func (rc *ReqClient) SetTlsCertFiles(certPemFilePath, keyPemFilePath string) *ReqClient {
	cert, err := tls.LoadX509KeyPair(certPemFilePath, keyPemFilePath)
	if err != nil {
		return nil
	}
	rc.tlsConf.Certificates = []tls.Certificate{cert}
	return rc
}
