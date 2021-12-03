package main

import (
	"context"
	"fmt"
	"log"

	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"github.com/qianxi0410/naive-rpc/config"
	"github.com/qianxi0410/naive-rpc/router"
	"github.com/qianxi0410/naive-rpc/transport"
)

func main() {
	y, err := config.NewYamlConfig("../conf/helloworld.yaml")
	if err != nil {
		log.Fatalf("load config error %s", err.Error())
	}

	addr := y.Read("server_addr", "")

	r := router.NewRouter()
	r.Forward("/hello", func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		pbreq := req.(*evangelion.Request)
		fmt.Printf("server recv req:%v", pbreq)

		pbrsp := &evangelion.Response{
			// Seqno:   pbreq.Seqno,
			// ErrCode: proto.Uint32(0),
			// ErrMsg:  proto.String("success"),
			Body: []byte("hello world"),
		}
		return pbrsp, nil
	})

	r.Forward("/ping", func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		pbreq := req.(*evangelion.Request)
		fmt.Printf("server recv req:%v", pbreq)

		pbrsp := &evangelion.Response{
			// Seqno:   pbreq.Seqno,
			// ErrCode: proto.Uint32(0),
			// ErrMsg:  proto.String("success"),
			Body: []byte("pong"),
		}
		return pbrsp, nil
	})

	tcpSvr, err := transport.NewTcpServerTransport(context.TODO(),
		"tcp4", addr, evangelion.NAME, transport.WithRouter(r))
	if err != nil {
		panic(err)
	}
	fmt.Printf("server is on %s\n", tcpSvr.Address())
	tcpSvr.ListenAndServe()
}
