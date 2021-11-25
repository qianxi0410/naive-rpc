package errors

import "fmt"

type errorType int

const (
	errorTypeNil = errorType(iota) // nil error
	Framework                      // framework error
	Business                       // business error
)

func (t errorType) String() string {
	switch t {
	case errorTypeNil:
		return "nil"
	case Framework:
		return "framework error"
	case Business:
		return "bussiness error"
	default:
		return "unknown error"
	}
}

// the struct of error
//
type Error struct {
	Code int       // error code
	Msg  string    // error msg
	Type errorType // error type
}

// New returns a new error
// frameworkd error should be defined in frameworkd
// there is no need to export
// users should only care about the type of error
func New(code int, msg string, errType errorType) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
		Type: errType,
	}
}

// impl the Go's internal `error` interface
// the %s in errorType will implict call errorType's String method
func (e *Error) Error() string {
	return fmt.Sprintf("err code: %d, error msg: %s, error type: %s", e.Code, e.Msg, e.Type)
}
