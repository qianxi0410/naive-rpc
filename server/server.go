package server

import (
	"context"
	"fmt"

	"github.com/qianxi0410/naive-rpc/router"
	"github.com/qianxi0410/naive-rpc/transport"
)

// Service represents a server instance (a server process),
// Any server can include more than one service, i.e, any server can be
// plugged into multile modules, like TcpServerTransport, UdpServerTransport, Broker, etc.
// By this way, we can implement more modules to extend server's abilities.
type Service struct {
	// service name
	name string
	// context
	ctx    context.Context
	cancel context.CancelFunc

	// a service have some tranport instance
	// may include udp tcp http
	// but now we only have tcp
	// []transport.Tranport
	trans transport.Transport

	// net type
	net string
	// codec
	codec string
}

// create a new server with option
func NewService(name string, net, codec string) *Service {
	s := &Service{
		name:  name,
		net:   net,
		codec: codec,
	}
	s.ctx, s.cancel = context.WithCancel(context.TODO())

	return s
}

// block func
func (r *Service) ListenAndServe(ctx context.Context, router *router.Router, addr string) error {
	var err error

	// transport options
	toptions := []transport.Option{}
	if router == nil {
		return fmt.Errorf("router is nil")
	}
	toptions = append(toptions, transport.WithRouter(router))

	if r.net == "tcp" || r.net == "tcp4" || r.net == "tcp6" {
		r.trans, err = transport.NewTcpServerTransport(ctx, r.net, addr, r.codec, toptions...)
		if err != nil {
			return err
		}

	}

	go r.trans.ListenAndServe()

	select {
	case <-ctx.Done():
		return nil
	case <-r.trans.Closed():
		return nil
	}

}
