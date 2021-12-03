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
			Body: []byte("pong1"),
		}, nil
	})

	err := s.ListenAndServe(context.TODO(), r, "localhost:8080")
	if err != nil {
		panic(err)
	}
}
