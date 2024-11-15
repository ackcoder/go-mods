package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
)

const (
	StringLetter = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberLetter = "0123456789"
)

// RandStr 生成指定长度的 随机(字母+数值)字符串
func RandStr[T Number](n T) string {
	resBytes := make([]byte, n)
	size := len(StringLetter)
	for i := range resBytes {
		resBytes[i] = StringLetter[rand.Intn(size)]
	}
	return string(resBytes)
}

// RandNumStr 生成指定长度的 随机(数值)字符串
func RandNumStr[T Number](n T) string {
	resBytes := make([]byte, n)
	size := len(NumberLetter)
	for i := range resBytes {
		resBytes[i] = NumberLetter[rand.Intn(size)]
	}
	return string(resBytes)
}

// Md5File 生成32位文件md5码
func Md5File(fs io.Reader) string {
	w := md5.New()
	if _, err := io.Copy(w, fs); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", w.Sum(nil))
}

// Md5Str 生成32位md5码
func Md5Str(str string) string {
	w := md5.New()
	if _, err := w.Write([]byte(str)); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", w.Sum(nil))
}

// Md5Str16 生成16位md5码
func Md5Str16(str string) string {
	md5Str := Md5Str(str)
	if md5Str == "" {
		return ""
	}
	return md5Str[8:24]
}
