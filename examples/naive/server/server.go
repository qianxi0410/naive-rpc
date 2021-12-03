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
	s := server.NewService("naive", "tcp", evangelion.NAME)

	err := s.ListenAndServe(context.TODO(), r, "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}

	r2 := router.NewRouter()
	r2.Forward("/ping", func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		return &evangelion.Response{
			Body: []byte("hello qianxi"),
		}, nil
	})
}
