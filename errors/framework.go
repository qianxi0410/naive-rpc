package errors

import "github.com/qianxi0410/naive-rpc/internal/errors"

// util func to create a framework type error
// no need to expose to user
func newFrameworkError(code int, msg string) *errors.Error {
	return errors.New(code, msg, errors.Framework)
}

var (
	// 100* reprsent server error
	// server context terminate
	ServerCtxDoneErr = newFrameworkError(1000, "server context done")
	// server side can's connect
	ServerNotInitErr = newFrameworkError(1002, "server not initialized")
	// no session found
	SessionNotExistErr = newFrameworkError(1003, "session not exist")
	// can't map func
	RouteNotFoundErr = newFrameworkError(1004, "route not found")

	// 200* reoresent codec error
	// decode error
	CodecDecodeErr = newFrameworkError(2000, "decode error")
	// encode error
	CodecEncodeErr = newFrameworkError(2001, "encode error")

	// 300* represent transport error
	// message read error
	CodecReadErr = newFrameworkError(3000, "read error package")
	// message read incomplete
	CodecReadIncompleteErr = newFrameworkError(3001, "read incomplete package")
	// message invalid
	CodecReadInvalid = newFrameworkError(3002, "read invalid package")
	// package too big
	CodecRead2BigErr = newFrameworkError(3004, "read to big package")

	// 400* represent connection pool error
	// pool too many
	ExceedPoolLimitErr = newFrameworkError(4000, "connection pool too many")
	// pool terminate
	PoolClosedErr = newFrameworkError(4001, "connection pool closed")
	// conn terminate
	ConnClosedErr = newFrameworkError(4002, "connection closed")
)
