package main

import (
	"context"
	"fmt"
	"time"

	"github.com/qianxi0410/naive-rpc/client"
	"github.com/qianxi0410/naive-rpc/client/selector"
	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/protobuf/proto"
)

func main() {
	c := client.NewClient("registry", clientv3.Config{
		Endpoints:   []string{"127.0.0.1:3000"},
		DialTimeout: 5 * time.Second,
	}, selector.ROUND_ROUBIN, client.WithCodec(evangelion.NAME), client.WithTransportType(client.TCP),
	)
	v, err := c.Invoke(context.TODO(), &evangelion.Request{
		Rpcname: proto.String("/ping"),
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("from server: %s", v.(*evangelion.Response).Body)
}
