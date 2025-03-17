package httpreq

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"time"
)

type ReqClient struct {
	domain  string        //请求目的域, 如"http://xxx.com"
	timeout time.Duration //请求超时, 默认10s
	tlsConf *tls.Config

	client *http.Client
}

func New(domain string) *ReqClient {
	return &ReqClient{
		domain:  domain,
		timeout: 10 * time.Second,
		tlsConf: &tls.Config{},
	}
}

func (rc *ReqClient) newClient() {
	rc.client = &http.Client{
		Timeout: rc.timeout,
		Transport: &http.Transport{
			TLSClientConfig: rc.tlsConf,
		},
	}
}

// SetTimeout 设置请求超时
func (rc *ReqClient) SetTimeout(second int) *ReqClient {
	rc.timeout = time.Duration(second) * time.Second
	if rc.client != nil {
		rc.newClient()
	}
	return rc
}

// SetTlsClientVerify 设置客户端Tls证书校验 (双向认证)
//   - {certPemFilePath} xxx.crt/cert.pem (publicKey.pem)
//   - {keyPemFilePath} xxx.key/key.pem (privateKey.pem)
func (rc *ReqClient) SetTlsClientVerify(certPemFilePath, keyPemFilePath string) *ReqClient {
	// // 非PEM结构密钥文件的处理
	// keyBytes, err := os.ReadFile(keyPemFilePath) //读取私钥文件
	// if err != nil {
	// 	panic(err)
	// }
	// block, rest := pem.Decode(keyBytes) //将字节流解码为pem结构
	// if len(rest) > 0 {
	// 	panic("pem解码失败")
	// }
	// der, err := x509.DecryptPEMBlock(block, []byte("pem私钥密码")) //pem解密
	// if err != nil {
	// 	panic(err)
	// }
	// key, err := x509.ParsePKCS1PrivateKey(der) //转为PKCS1 RSA私钥
	// if err != nil {
	// 	panic(err)
	// }
	// keyPemBlock := pem.EncodeToMemory(&pem.Block{ //编码为新的pem结构私钥
	// 	Type: "RSA PRIVATE KEY",
	// 	Bytes: x509.MarshalPKCS1PrivateKey(key),
	// })
	// certPemBlock, err := os.ReadFile(certPemFilePath) //读取公钥文件
	// if err != nil {
	// 	panic(err)
	// }
	// cert, err := tls.X509KeyPair(certPemBlock, keyPemBlock)

	cert, err := tls.LoadX509KeyPair(certPemFilePath, keyPemFilePath)
	if err != nil {
		panic(err)
	}
	rc.tlsConf.Certificates = []tls.Certificate{cert}
	if rc.client != nil {
		rc.newClient()
	}
	return rc
}

// SetTlsServerSkipVerify 设置服务端Tls证书跳过校验
func (rc *ReqClient) SetTlsServerSkipVerify() *ReqClient {
	rc.tlsConf.InsecureSkipVerify = true
	if rc.client != nil {
		rc.newClient()
	}
	return rc
}

// SetTlsServerVerify 设置服务端Tls证书校验 (自签证书校验)
//   - {caCrtFilePath} xxx.crt/ca.crt
func (rc *ReqClient) SetTlsServerVerify(caCrtFilePath string) *ReqClient {
	caCert, err := os.ReadFile(caCrtFilePath)
	if err != nil {
		panic(err)
	}
	rootCAs := x509.NewCertPool()
	if ok := rootCAs.AppendCertsFromPEM(caCert); !ok {
		panic("无法加载 CA 证书")
	}
	rc.tlsConf.RootCAs = rootCAs
	if rc.client != nil {
		rc.newClient()
	}
	return rc
}
