package main

import (
	"context"
	"fmt"

	"github.com/qianxi0410/naive-rpc/client"
	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"google.golang.org/protobuf/proto"
)

func main() {
	c := client.NewClient("naive", client.WithAddr("127.0.0.1:8888"), client.WithCodec(evangelion.NAME), client.WithTransportType(client.TCP4))
	v, err := c.Invoke(context.TODO(), &evangelion.Request{
		Rpcname: proto.String("/ping"),
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("receive from server : %s\n", v.(*evangelion.Response).Body)
}
