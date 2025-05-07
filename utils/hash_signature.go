package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type sign struct{}

// Sign 签名工具方法
//
//   - 关于 HMAC
//
// HMAC (Hash-based Message Authentication Code, 基于哈希的消息认证码)
// 用于消息身份验证、消息完整性校验、密钥派生等，因为单纯使用哈希函数无法认证来源与防止篡改
//
//   - 关于哈希算法
//
// SHA 哈希算法是散列算法的一种、主流有 SHA-1, SHA-2(SHA-256, SHA-384, SHA-512), SM3 实现
//
//   - 关于编码方式
//
// Hex 十六进制编码
// 每字节对应两字符、占用约 200% 空间，用于颜色代码/mac地址等
//
// Base64 编码
// 每字节对应三字符、占用约 133% 空间、结尾可能有一或二个"="填充，用于邮件附件/url参数等
var Sign = sign{}

// HmacSha256Hex 生成签名
//   - {data} 待签名数据
//   - {secret} 密钥
//
// 注: 更多详情请阅读 Sign 实例注释
func (s sign) HmacSha256Hex(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// HmacSha256B64 生成签名
//   - {data} 待签名数据
//   - {secret} 密钥
//
// 注: 更多详情请阅读 Sign 实例注释
func (s sign) HmacSha256B64(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HmacSha1Hex 生成签名
//   - {data} 待签名数据
//   - {secret} 密钥
//
// 注: 更多详情请阅读 Sign 实例注释
func (s sign) HmacSha1Hex(data, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// HmacSha1B64 生成签名
//   - {data} 待签名数据
//   - {secret} 密钥
//
// 注: 更多详情请阅读 Sign 实例注释
func (s sign) HmacSha1B64(data, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
