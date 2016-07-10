package Errors
/*
自定义的errors类型，支持格式化输出错误字符串
*/

//package errors

import (
	"fmt"
)

func New(format string, arg ...interface{}) error {
	s := fmt.Sprintf(format, arg...)
	return &myError{s}
}

type myError struct {
	s string
}

func (e *myError) Error() string {
	return e.s
}
