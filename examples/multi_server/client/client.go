package main

import (
	"context"
	"fmt"

	"github.com/qianxi0410/naive-rpc/client"
	"github.com/qianxi0410/naive-rpc/client/selector"
	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"google.golang.org/protobuf/proto"
)

func main() {
	c := client.NewClient("multi", client.WithCodec(evangelion.NAME), client.WithTransportType(client.TCP),
		client.WithSelector(selector.NewIPSelector("tcp", []string{"127.0.0.1:9090", "127.0.0.1:8080"})))
	v, err := c.Invoke(context.TODO(), &evangelion.Request{
		Rpcname: proto.String("/ping"),
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("from server: %s", v.(*evangelion.Response).Body)
}
