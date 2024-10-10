package utils

import (
	"strings"
)

// Must 必须不为Nil
func Must(in any, msg ...string) {
	if inE, ok := in.(error); ok {
		msg = append(msg, inE.Error())
	}
	panicIfCondition(in != nil, msg...)
}

// MustNil 必须为Nil
func MustNil(in any, msg ...string) {
	if inE, ok := in.(error); ok {
		msg = append(msg, inE.Error())
	}
	panicIfCondition(in == nil, msg...)
}

// MustZero 必须为0
func MustZero(in int, msg ...string) {
	panicIfCondition(in == 0, msg...)
}

// MustZeroN 必须不为0
func MustZeroN(in int, msg ...string) {
	panicIfCondition(in != 0, msg...)
}

// MustTrue 必须为true
func MustTrue(in bool, msg ...string) {
	panicIfCondition(in, msg...)
}

// MustFalse 必须为false
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
