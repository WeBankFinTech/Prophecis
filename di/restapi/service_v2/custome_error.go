package service_v2

import "fmt"

type InputParamsError struct {
	Field  string
	Reason string
}

func (e InputParamsError) Error() string {
	if e.Reason == "" {
		return fmt.Sprintf("输入%s字段有误", e.Field)
	} else {
		return fmt.Sprintf("输入%s字段有误: %s", e.Field, e.Reason)
	}
}

type InternalError struct {
}

func (e InternalError) Error() string {
	return "内部错误"
}

type AccessError struct {
}

func (e AccessError) Error() string {
	return "权限错误"
}
