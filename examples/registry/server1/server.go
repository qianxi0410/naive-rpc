package main

import (
	"context"
	"time"

	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"github.com/qianxi0410/naive-rpc/router"
	"github.com/qianxi0410/naive-rpc/server"
	"go.etcd.io/etcd/clientv3"
)

func main() {
	s := server.NewService("registry", "tcp", "localhost:8080", evangelion.NAME, clientv3.Config{
		Endpoints:   []string{"127.0.0.1:3000"},
		DialTimeout: 5 * time.Second,
	})

	r := router.NewRouter()
	r.Forward("/ping", func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		return &evangelion.Response{
			Body: []byte("pong1"),
		}, nil
	})

	err := s.ListenAndServe(context.TODO(), r)
	if err != nil {
		panic(err)
	}
}
