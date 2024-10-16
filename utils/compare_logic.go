package utils

import (
	"strings"
)

type integerNumber interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// Must 必须无错误 否则Panic
func Must(in error, msg ...string) {
	msg = append(msg, in.Error())
	panicIfCondition(in != nil, msg...)
}

// MustZero 必须为0 否则Panic
func MustZero[T integerNumber](in T, msg ...string) {
	panicIfCondition(in == 0, msg...)
}

// MustZeroN 必须不为0 否则Panic
func MustZeroN[T integerNumber](in T, msg ...string) {
	panicIfCondition(in != 0, msg...)
}

// MustTrue 必须为True 否则Panic
func MustTrue(in bool, msg ...string) {
	panicIfCondition(in, msg...)
}

// MustFalse 必须为False 否则Panic
func MustFalse(in bool, msg ...string) {
	panicIfCondition(!in, msg...)
}

// Mustxxx 公共逻辑
func panicIfCondition(condition bool, msg ...string) {
	if !condition {
		return
	}
	if len(msg) != 0 {
		panic(strings.Join(msg, " "))
	} else {
		panic("Mustxxx compare panic")
	}
}
