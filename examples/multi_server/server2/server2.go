package main

import (
	"context"

	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"github.com/qianxi0410/naive-rpc/router"
	"github.com/qianxi0410/naive-rpc/server"
)

func main() {
	s := server.NewService("multi", "tcp", evangelion.NAME)

	r := router.NewRouter()
	r.Forward("/ping", func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		return &evangelion.Response{
			Body: []byte("pong2"),
		}, nil
	})

	err := s.ListenAndServe(context.TODO(), r, "localhost:9090")
	if err != nil {
		panic(err)
	}
}
