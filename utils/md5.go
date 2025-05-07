package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

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
