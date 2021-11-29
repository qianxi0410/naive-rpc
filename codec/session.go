package codec

import (
	"context"
	"sync"
)

// define the rpc context for a request
// the session init instally in client side
// when server side receive the requset, session init in server side
// maybe call `RpcContext`
type Session interface {
	RPCName() string

	Request() interface{}

	Response() interface{}

	SetResponse(rsp interface{})

	SetError(err error)
}

// BaseSession implements some basic methods defined in `Session`
type BaseSession struct {
	ReqHead interface{}
	RspHead interface{}
}

func (r *BaseSession) Request() interface{} {
	if r != nil {
		return r.ReqHead
	}
	return nil
}

func (r *BaseSession) Response() interface{} {
	if r != nil {
		return r.RspHead
	}
	return nil
}

func (r *BaseSession) SetResponse(rsp interface{}) {
	if r != nil {
		r.RspHead = rsp
	}
}

var (
	sessionMutex sync.RWMutex
	builders     = map[string]SessionBuilder{}
)

// SessionBuilder when extending protocols, SessionBuilder should be
// implemented and registered to help build the `Session`.
type SessionBuilder interface {
	Build(reqHead interface{}) (Session, error)
}

// RegisterSessionBuilder register extended SessionBuilder for protocol `proto`
func RegisterSessionBuilder(proto string, builder SessionBuilder) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	builders[proto] = builder
}

// GetSessionBuilder return SessionBuilder for protocol `proto`
func GetSessionBuilder(proto string) SessionBuilder {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	return builders[proto]
}

const sessionKey = "qianxi"

// SessionFromContext return Session carried by `ctx`
func SessionFromContext(ctx context.Context) Session {
	v := ctx.Value(sessionKey)
	session, ok := v.(Session)
	if !ok {
		return nil
	}
	return session
}

// ContextWithSession return new context carrying value `session`
func ContextWithSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}
