package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/qianxi0410/naive-rpc/client/pool"
	"github.com/qianxi0410/naive-rpc/client/selector"
	"github.com/qianxi0410/naive-rpc/client/transport"
	"github.com/qianxi0410/naive-rpc/codec"
	"github.com/qianxi0410/naive-rpc/codec/evangelion"
)

var defaultPoolFactory = pool.NewConnPoolFactory(
	pool.WithMinIdle(2),
	pool.WithMaxIdle(4),
	pool.WithMaxActive(8),
	pool.WithDialTimeout(time.Second*2),
	pool.WithIdleTimeout(time.Minute*5),
	pool.WithMaxConnLifetime(time.Minute*30),
	pool.WithWait(true),
)

type Client interface {
	Invoke(ctx context.Context, reqHead interface{}, opts ...Option) (rspHead interface{}, err error)
}

type client struct {
	// srv name
	Name          string
	Addr          string
	Codec         codec.Codec
	Selector      selector.Selector
	Transport     transport.Transport
	TransportType TransportType
	RpcType       RpcType
}

func (r *client) Invoke(ctx context.Context, reqHead interface{}, opts ...Option) (rspHead interface{}, err error) {

	var (
		network string
		address string
	)

	if r.Addr != "" && r.TransportType.Valid() {
		network = r.TransportType.String()
		fmt.Println(network)
		address = strings.TrimPrefix(r.Addr, "ip://")
	} else if r.Name != "" && r.Selector != nil {
		node, err := r.Selector.Select(r.Name)
		if err != nil {
			return nil, err
		}
		network = node.Network
		address = node.Address
	}

	rsp, err := r.Transport.Send(ctx, network, address, reqHead)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func NewClient(name string, opts ...Option) Client {

	c := &client{
		Name:          name,
		Selector:      nil,
		TransportType: TCP,
		Transport:     &transport.TcpTransport{},
		// Address:      addr,
		Codec:   codec.ClientCodec(evangelion.NAME),
		RpcType: SendRecv,
	}

	for _, o := range opts {
		o(c)
	}
	return c
}
