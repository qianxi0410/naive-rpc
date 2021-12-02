package main

import (
	"context"
	"fmt"

	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"github.com/qianxi0410/naive-rpc/router"
	"github.com/qianxi0410/naive-rpc/transport"
	"google.golang.org/protobuf/proto"
)

func main() {
	r := router.NewRouter()
	r.Forward("/hello", func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		pbreq := req.(*evangelion.Request)
		fmt.Printf("server recv req:%v", pbreq)

		pbrsp := &evangelion.Response{
			Seqno:   pbreq.Seqno,
			ErrCode: proto.Uint32(0),
			ErrMsg:  proto.String("success"),
			Body:    []byte("hello " + *pbreq.Userid),
		}
		return pbrsp, nil
	})

	tcpSvr, err := transport.NewTcpServerTransport(context.TODO(),
		"tcp4", "127.0.0.1:8888", evangelion.NAME, transport.WithRouter(r))
	if err != nil {
		panic(err)
	}
	fmt.Printf("server is on %s", tcpSvr.Address())
	tcpSvr.ListenAndServe()
}
