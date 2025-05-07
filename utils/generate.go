package utils

import "math/rand/v2"

const (
	StringLetter = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberLetter = "0123456789"
)

// RandStr 生成指定长度的 随机(字母+数值)字符串
func RandStr[T Number](n T) string {
	resBytes := make([]byte, n)
	size := len(StringLetter)
	for i := range resBytes {
		resBytes[i] = StringLetter[rand.IntN(size)]
	}
	return string(resBytes)
}

// RandNumStr 生成指定长度的 随机(数值)字符串
func RandNumStr[T Number](n T) string {
	resBytes := make([]byte, n)
	size := len(NumberLetter)
	for i := range resBytes {
		resBytes[i] = NumberLetter[rand.IntN(size)]
	}
	return string(resBytes)
}
