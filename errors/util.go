package errors

import "github.com/qianxi0410/naive-rpc/internal/errors"

// return a business error
// expose this method to user
func New(code int, msg string) *errors.Error {
	return errors.New(code, msg, errors.Business)
}

// to judge the error if is framework error
// is nil -> false
// is not *Error -> false
// is not Framework type -> false
func IsFrameworkError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*errors.Error)
	if !ok {
		return false
	}

	return e.Type == errors.Framework
}

// to judge the error if is business error
// is nil -> false
// is not *Error -> false
// is not Business type -> false
func IsBusinessError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*errors.Error)
	if !ok {
		return false
	}

	return e.Type == errors.Business
}

// reutrn code of error if error type is *Error
// else reuturn 0
func ErrorCode(err error) int {
	if err == nil {
		return 0
	}

	e, ok := err.(*errors.Error)
	if !ok {
		return 0
	}
	return e.Code
}

// return msg of error
func ErrorMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
