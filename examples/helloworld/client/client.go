package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qianxi0410/naive-rpc/client"
	"github.com/qianxi0410/naive-rpc/client/selector"
	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"github.com/qianxi0410/naive-rpc/config"
	"google.golang.org/protobuf/proto"
)

func main() {
	y, err := config.NewYamlConfig("../conf/helloworld.yaml")
	if err != nil {
		log.Fatalf("read config file err: %s", err.Error())
	}

	server_addr := y.Read("server_addr", "")
	// user_id := y.Read("user_id", "")
	// user_key := y.Read("user_key", "")

	cli := client.NewClient("test",
		client.WithAddr(server_addr),
		client.WithCodec(evangelion.NAME),
		client.WithTransportType(client.TCP4),
		client.WithSelector(selector.NewIPSelector("tcp4", []string{server_addr})),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	reqHead := &evangelion.Request{
		// Seqno:   proto.Uint64(1000),
		// Appid:   proto.String("strong"),
		Rpcname: proto.String("/ping"),
		// Userid:  proto.String(user_id),
		// Userkey: proto.String(user_key),
		// Version: proto.Uint32(0),
	}
	fmt.Println(reqHead)

	v, err := cli.Invoke(ctx, reqHead)
	if err != nil {
		fmt.Println(err)
		return
	}

	rspHead := v.(*evangelion.Response)
	fmt.Println("receive from server : " + string(rspHead.GetBody()))
	// fmt.Println("client recv response:", rspHead)
}
