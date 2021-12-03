package main

import (
	"context"

	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"github.com/qianxi0410/naive-rpc/router"
	"github.com/qianxi0410/naive-rpc/server"
)

func main() {
	r := router.NewRouter()

	r.Forward("/ping", func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		return &evangelion.Response{
			Body: []byte("hello world"),
		}, nil
	})

	// naiverpc.ListenAndServe(r, naiverpc.WithConf("../conf/service.yaml"))
	s := server.NewService("naive")

	err := s.ListenAndServe(context.TODO(), "tcp", "127.0.0.1:8888", evangelion.NAME, server.WithRouter(r))
	if err != nil {
		panic(err)
	}
}
