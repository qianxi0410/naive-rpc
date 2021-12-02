package evangelion

import (
	"fmt"

	"github.com/qianxi0410/naive-rpc/codec"
	"github.com/qianxi0410/naive-rpc/errors"
	"google.golang.org/protobuf/proto"
)

type DefaultSession struct {
	codec.BaseSession
}

func (r *DefaultSession) RPCName() string {
	return r.Request().(*Request).GetRpcname()
}

func (r *DefaultSession) SetError(err error) {
	var code int
	var msg string
	rsp := r.Response().(Response)

	if errors.IsFrameworkError(err) {
		code = errors.ErrorCode(err)
		msg = errors.ErrorMsg(err)
	} else {
		code = 10000
		msg = err.Error()
	}
	rsp.ErrCode = proto.Uint32(uint32(code))
	rsp.ErrMsg = proto.String(msg)
}

type DefaultSessionBuilder struct{}

func (r *DefaultSessionBuilder) Build(reqHead interface{}) (codec.Session, error) {
	req, ok := reqHead.(*Request)
	if !ok {
		return nil, fmt.Errorf("req:%v not *evangelion.Request", req)
	}

	rspHead := &Response{}
	rspHead.Seqno = req.Seqno

	session := &DefaultSession{
		codec.BaseSession{
			ReqHead: req,
			RspHead: rspHead,
		},
	}

	return session, nil
}
