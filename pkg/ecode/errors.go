package ecode

import (
	"errors"
	"fmt"
)

const (
	// UnknownReason 未知错误
	UnknownReason = "unknown"
	// UnknownReasonCode 未知错误代码
	UnknownReasonCode = 500
)

var _ error = &StatusError{}

// StatusError ..
type StatusError struct {
	Code       int           `json:"code"`
	Message    string        `json:"message"`
	Debug      []interface{} `json:"debug"`
	*ErrorInfo `json:"err"`
}

// Error ..
func (e *StatusError) Error() string {
	return fmt.Sprintf("errs: code = %d err = %s errDetails = %+v", e.Code, e.ErrorInfo, e.Debug)
}

// WithDebugs 记录debug相关的信息,入参出参
func (e *StatusError) WithDebugs(debugs ...interface{}) *StatusError {
	e.Debug = append(e.Debug, debugs...)
	return e
}

// Is ..
func (e *StatusError) Is(target error) bool {
	err, ok := target.(*StatusError)
	if ok {
		return e.Code == err.Code
	}
	return false
}

// ErrorInfo 详情信息
type ErrorInfo struct {
	// 子类枚举错误
	Reason string `json:"reason"`
	// 前端toast
	Toast string `json:"msg"`
	// 备注等其他信息
	Detail interface{} `json:"detail"`
}

// String ..
func (e *ErrorInfo) String() string {
	return fmt.Sprintf("(reason: %s toast: %s detail: %s)", e.Reason, e.Toast, e.Detail)
}

// Reason 查询err底层数据结构
func Reason(err error) *ErrorInfo {
	if se := new(StatusError); errors.As(err, &se) {
		if se.ErrorInfo != nil {
			return se.ErrorInfo
		}
	}
	return &ErrorInfo{Reason: UnknownReason}
}

// HTTPCodeAndReason 返回HTTP-code&响应值
func HTTPCodeAndReason(err error) (int, *ErrorInfo) {
	se := Unknown(UnknownReason, err.Error(), "")
	if !errors.As(err, &se) {
		return UnknownReasonCode, Reason(se)
	}
	if code, exists := gRPC2HTTPCode[se.Code]; exists {
		return code, se.ErrorInfo
	}
	return UnknownReasonCode, se.ErrorInfo
}

// Error 实例化错误信息
func Error(code int, err *ErrorInfo, debugs ...interface{}) error {
	return &StatusError{Code: code, ErrorInfo: err, Debug: debugs}
}

// Errorf ..
func Errorf(code int, err *ErrorInfo, format string, a ...interface{}) error {
	return Error(code, err, fmt.Sprintf(format, a...))
}
