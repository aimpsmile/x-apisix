package ecode

import (
	"github.com/pkg/errors"
)

// 生成新的错误code码
func New(code int32, msg string) Codes {
	return &Code{code: code, msg: msg}
}

// Codes ecode error interface which has a code & message.
type Codes interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int32
	//	show msg to client
	Msg() string
}

// A Code is an int error code spec.
type Code struct {
	code int32
	msg  string
}

func (e *Code) Error() string {
	return e.msg
}

// Code return error code
func (e Code) Code() int32 { return e.code }
func (e Code) Msg() string { return e.msg }

//	error types are converted by Code
func Error(code Codes) error {
	return code
}

// Cause cause from error to ecode.
func Cause(e error) Codes {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Codes)
	if ok {
		return ec
	}
	return UnknownErr
}

// Equal equal a and b by code int.
func Equal(a, b Codes) bool {
	if a == nil {
		a = OK
	}
	if b == nil {
		b = OK
	}
	return a.Code() == b.Code()
}

// EqualError code equal code
func EqualError(code Codes, err error) bool {
	return Cause(err).Code() == code.Code()
}
