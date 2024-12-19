package xerrors

import (
	"errors"
	"fmt"
)

type Error struct {
	Code   int
	Msg    string
	parent error
}

func (e *Error) Error() string {
	return fmt.Sprintf("code=%d msg=%s", e.code, e.msg)
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

func NewParamsError(err error) error {
	return &Error{
		Code:   400,
		parent: err,
		Msg:    "参数错误",
	}
}
