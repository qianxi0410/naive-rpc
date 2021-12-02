package main

import (
	"context"
	"fmt"
	"time"

	"github.com/qianxi0410/naive-rpc/client"
	"github.com/qianxi0410/naive-rpc/client/selector"
	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"google.golang.org/protobuf/proto"
)

func main() {
	cli := client.NewClient("test",
		client.WithAddr("127.0.0.1:8888"),
		client.WithCodec(evangelion.NAME),
		client.WithTransportType(client.TCP4),
		client.WithSelector(selector.NewIPSelector("tcp4", []string{"127.0.0.1:8888"})),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	reqHead := &evangelion.Request{
		Seqno:   proto.Uint64(1000),
		Appid:   proto.String("strong"),
		Rpcname: proto.String("/hello"),
		Userid:  proto.String("qianxi"),
		Userkey: proto.String("cat"),
		Version: proto.Uint32(0),
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
