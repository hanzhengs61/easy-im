package errorx

import (
	"errors"
	"fmt"
)

// AppError 业务错误，携带错误码和用户可见提示
// 实现了 error 接口，可以在任何地方像普通 error 一样使用
type AppError struct {
	Code    int    // 业务错误码
	Msg     string // 用户可见提示（可被业务层覆盖）
	Details string // 内部详情，只写日志，不返回给前端
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("code=%d, msg=%s, details=%s", e.Code, e.Msg, e.Details)
	}
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Msg)
}

// New 创建业务错误，使用错误码默认提示
func New(code int) *AppError {
	return &AppError{
		Code: code,
		Msg:  Message(code),
	}
}

// NewWithMsg 创建业务错误，自定义用户提示
func NewWithMsg(code int, msg string) *AppError {
	return &AppError{
		Code: code,
		Msg:  msg,
	}
}

// NewWithDetails 创建业务错误，附带内部调试详情（不暴露给用户）
func NewWithDetails(code int, details string) *AppError {
	return &AppError{
		Code:    code,
		Msg:     Message(code),
		Details: details,
	}
}

// Wrap 包装底层 error 为业务错误（用于 dao 层错误向上传递）
func Wrap(code int, err error) *AppError {
	details := ""
	if err != nil {
		details = err.Error()
	}
	return &AppError{
		Code:    code,
		Msg:     Message(code),
		Details: details,
	}
}

// IsAppError 判断 error 是否为业务错误
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
