package xerrors

import (
	"errors"
	"fmt"
)

// 错误码+维护状态设置
// 200 = 正常
// 400 = 参数传递错误, 前端自行toast
// 401 = 业务逻辑错误，前端可以弹出后端的msg
// 500 = 服务器内部错误，前端弹出服务器内部错误

type Error struct {
	Code   int
	Msg    string
	parent error
}

func (e *Error) Error() string {
	str := fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
	if e.parent != nil {
		str += " parent=" + e.parent.Error()
	}
	return str
}

func (e *Error) Unwrap() error {
	return e.parent
}

func IsError(err error) (*Error, bool) {
	var myError *Error
	if errors.As(err, &myError) {
		return myError, true
	}
	return nil, false
}

// NewParamsError 前端提示参数错误，入参是err
func NewParamsError(err error) error {
	return &Error{
		Code:   400,
		parent: err,
		Msg:    "参数错误",
	}
}

// NewParamsErrors 前端提示参数错误，入参是 string
func NewParamsErrors(s string) error {
	return &Error{
		Code: 400,
		Msg:  s,
	}
}

// NewCustomError 前端直接弹出后端的提示语
func NewCustomError(s string) error {
	return &Error{
		Code: 401,
		Msg:  s,
	}
}
