package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

type pwd struct {
	// scrypt params
	s_N       int //迭代次数
	s_r       int //块大小
	s_p       int //并行度
	s_keyLen  int //生成密钥长度
	s_saltLen int //随机盐值长度

	// argon2 params
	a_t       uint32 //迭代数
	a_m       uint32 //使用内存大小
	a_p       uint8  //并行线程数
	a_keyLen  uint32 //生成密钥长度
	a_saltLen int    //随机盐值长度
}

// Pwd 密码哈希工具
//
//   - 关于 KDF (Key Derivation Functions, 密钥派生函数)
//
// 本包使用的 KDF 实现有: bcrypt, scrypt, argon2
//
// bcrypt
// 最早出现, 通过加盐(salt)防止彩虹表攻击, 配置工作因子(cost)调节计算时间以对抗硬件加速暴力破解
//
// scrypt
// 晚于 bcrypt 出现, 通过内存密集计算来抵抗GPU/ASIC/FPGA等硬件加速暴力破解
//
// argon2
// 密码哈希竞赛胜者, 在配置合理情况下, 比 bcrypt,scrypt 有更强的抗破解性
// 有三个变种实现 argon2d, argon2i, argon2id, 本包使用 argon2id
var Pwd = pwd{
	s_N:       32768,
	s_r:       8,
	s_p:       1,
	s_keyLen:  32,
	s_saltLen: 32,

	a_t:       1,
	a_m:       64 * 1024,
	a_p:       4,
	a_keyLen:  32,
	a_saltLen: 32,
}

// MakeByBcrypt 生成密码哈希
//   - {pwdTxt} 用户输入的密码明文
//
// 注: 更多详情请阅读 Pwd 实例注释
func (p pwd) MakeByBcrypt(pwdTxt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwdTxt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckByBcrypt 验证密码哈希
//   - {pwdTxt} 用户输入的密码明文
//   - {hash} 存储的密码哈希
//
// 注: 更多详情请阅读 Pwd 实例注释
func (p pwd) CheckByBcrypt(pwdTxt, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwdTxt))
	return err == nil
}

// ============================================================

// EditScryptParam 调整密码哈希参数
//   - {sN} 迭代次数, 2017年推荐交互式登录参考值为 32768 (2^15)
//   - {sr} 块大小, 2017年推荐交互式登录参考值为 8
//   - {sp} 并行度, 2017年推荐交互式登录参考值为 1
//   - {skeyLen} 生成密钥长度, 不影响计算时内存与时间消耗, 一般为 2 的幂次方
//   - {ssaltLen} 随机盐值长度, 仅用于生成和校验时
//
// 所需内存 = 128 * r * N * p
//
// 注: 更多详情请阅读 Pwd 实例注释
func (p pwd) EditScryptParam(sN, sr, sp, skeyLen, ssaltLen int) {
	p.s_N = sN
	p.s_r = sr
	p.s_p = sp
	p.s_keyLen = skeyLen
	p.s_saltLen = ssaltLen
}

// MakeByScrypt 生成密码哈希
//   - {pwdTxt} 用户输入的密码明文
//
// @return string Base64编码的密码哈希, 实际是盐值+密码哈希
//
// 建议在不同机器上通过比对生成结果的时间, 用 EditScryptParam 确定最佳参数
//
// 注: 更多详情请阅读 Pwd 实例注释
func (p pwd) MakeByScrypt(pwdTxt string) (string, error) {
	// 注: 随机盐值. 因在密码生成领域, 使用的是比 math/rand 更昂贵的密码学随机数包 crypto/rand
	saltBytes := make([]byte, p.s_saltLen)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}
	pwdHash, err := scrypt.Key([]byte(pwdTxt), saltBytes, p.s_N, p.s_r, p.s_p, p.s_keyLen)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(append(saltBytes, pwdHash...)), nil
}

// CheckByScrypt 验证密码哈希
//   - {pwdTxt} 用户输入的密码明文
//   - {hash} 存储的密码哈希
//
// 注: 更多详情请阅读 Pwd 实例注释
func (p pwd) CheckByScrypt(pwdTxt, hash string) bool {
	hashBytes, _ := base64.StdEncoding.DecodeString(hash)

	pwdHash, err := scrypt.Key([]byte(pwdTxt), hashBytes[:p.s_saltLen], p.s_N, p.s_r, p.s_p, p.s_keyLen)
	if err != nil {
		return false
	}
	// 不使用 string(pwdHash) == hash 比较
	// 用 subtle.ConstantTimeCompare() 避免时序攻击
	return subtle.ConstantTimeCompare(pwdHash, hashBytes[p.s_saltLen:]) == 1
}

// ============================================================

// EditArgon2Param 调整密码哈希参数
//   - {at} 迭代数, 根据 RFC 草案建议合理值为 1
//   - {am} 使用内存大小, 根据 RFC 草案建议合理值为 64 * 1024
//   - {ap} 并行线程数, 默认 4
//   - {akeyLen} 生成密钥长度, 默认 32
//   - {asaltLen} 随机盐值长度, 仅用于生成和校验时
//
// 注: 更多详情请阅读 Pwd 实例注释
func (p pwd) EditArgon2Param(at, am uint32, ap uint8, akeyLen uint32, asaltLen int) {
	p.a_t = at
	p.a_m = am
	p.a_p = ap
	p.a_keyLen = akeyLen
	p.a_saltLen = asaltLen
}

// MakeByArgon2 生成密码哈希
//   - {pwdTxt} 用户输入的密码明文
//
// @return string Base64编码的密码哈希, 实际是盐值+密码哈希
//
// 注: 更多详情请阅读 Pwd 实例注释
func (p pwd) MakeByArgon2(pwdTxt string) (string, error) {
	// 注: 随机盐值. 因在密码生成领域, 使用的是比 math/rand 更昂贵的密码学随机数包 crypto/rand
	saltBytes := make([]byte, p.a_saltLen)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", err
	}
	pwdHash := argon2.IDKey([]byte(pwdTxt), saltBytes, p.a_t, p.a_m, p.a_p, p.a_keyLen)
	return base64.StdEncoding.EncodeToString(append(saltBytes, pwdHash...)), nil
}

// CheckByArgon2 验证密码哈希
//   - {pwdTxt} 用户输入的密码明文
//   - {hash} 存储的密码哈希
//
// 注: 更多详情请阅读 Pwd 实例注释
func (p pwd) CheckByArgon2(pwdTxt, hash string) bool {
	hashBytes, _ := base64.StdEncoding.DecodeString(hash)

	pwdHash := argon2.IDKey([]byte(pwdTxt), hashBytes[:p.a_saltLen], p.a_t, p.a_m, p.a_p, p.a_keyLen)
	// 不使用 string(pwdHash) == hash 比较
	// 用 subtle.ConstantTimeCompare() 避免时序攻击
	return subtle.ConstantTimeCompare(pwdHash, hashBytes[p.a_saltLen:]) == 1
}
