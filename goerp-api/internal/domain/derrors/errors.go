package derrors

import (
	"errors"
	"fmt"
)

type DomainError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *DomainError) Error() string {
	return e.Message
}

func New(code int, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrInvalidParam        = New(400001, "参数无效")
	ErrUserNotFound        = New(404001, "用户不存在")
	ErrInvalidCredentials  = New(401001, "用户名或密码错误")
	ErrVerificationExpired = New(401002, "验证码已过期或无效")
	ErrInvalidVerification = New(401003, "验证码错误")
	ErrInternalError       = New(500001, "服务器内部错误")
)

func FromError(err error) *DomainError {
	var dErr *DomainError
	if errors.As(err, &dErr) {
		return dErr
	}
	return New(500000, err.Error())
}

func (e *DomainError) WithMessage(msg string) *DomainError {
	return New(e.Code, fmt.Sprintf("%s: %s", e.Message, msg))
}
