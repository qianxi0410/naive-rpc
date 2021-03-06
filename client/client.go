package client

import (
	"context"
	"strings"
	"time"

	"github.com/kpango/glg"
	"github.com/qianxi0410/naive-rpc/client/pool"
	"github.com/qianxi0410/naive-rpc/client/selector"
	"github.com/qianxi0410/naive-rpc/client/transport"
	"github.com/qianxi0410/naive-rpc/codec"
	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"go.etcd.io/etcd/clientv3"
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
		address = strings.TrimPrefix(r.Addr, "ip://")
		// FIXME:
	} else if r.Name != "" && r.Selector != nil {
		node, err := r.Selector.Select(r.Name)
		if err != nil {
			return nil, err
		}

		network = node.Network
		address = node.Address
		// addrs, err := r.registry.GetAddrs(r.Name, r.TransportType.String())
		// if err != nil || len(addrs) == 0 {
		// 	return nil, err
		// }
		// node, err := r.Selector.Select(r.Name)
		// if err != nil {
		// 	return nil, err
		// }
		// // network = r.TransportType.String()
		// // address = addrs[0]
		// network = node.Network
		// address = node.Address
	}

	rsp, err := r.Transport.Send(ctx, network, address, reqHead)
	if err != nil {
		return nil, err
	}
	glg.Successf("client invoke remote method success")
	return rsp, nil
}

func NewClient(name string, conf clientv3.Config, typ selector.SelectorType, opts ...Option) Client {

	c := &client{
		Name:          name,
		Selector:      nil,
		TransportType: TCP,
		Transport:     &transport.TcpTransport{},
		Codec:         codec.ClientCodec(evangelion.NAME),
		RpcType:       SendRecv,
	}

	for _, o := range opts {
		o(c)
	}

	c.Selector = selector.NewIPSelector(c.Name, c.TransportType.String(), typ, conf)
	glg.Successf("%s clinet is init", name)
	return c
}
