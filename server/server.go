package server

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/qianxi0410/naive-rpc/registry"
	"github.com/qianxi0410/naive-rpc/router"
	"github.com/qianxi0410/naive-rpc/transport"
	"go.etcd.io/etcd/clientv3"
)

// Service represents a server instance (a server process),
// Any server can include more than one service, i.e, any server can be
// plugged into multile modules, like TcpServerTransport, UdpServerTransport, Broker, etc.
// By this way, we can implement more modules to extend server's abilities.
type Service struct {
	// service name
	Name string
	// context
	ctx    context.Context
	cancel context.CancelFunc

	// a service have some tranport instance
	// may include udp tcp http
	// but now we only have tcp
	// []transport.Tranport
	trans transport.Transport

	// net type
	Net string
	// addr type
	Addr string
	// codec
	codec string
	// for etcd, like an uniqie id
	Id string
	// etcd
	registry registry.Registry
}

// create a new server with option
func NewService(name string, net, addr, codec string, c clientv3.Config) *Service {
	s := &Service{
		Name:     name,
		Net:      net,
		Addr:     addr,
		codec:    codec,
		Id:       string(md5.New().Sum([]byte(name + addr))),
		registry: registry.NewEvaRegistry(c),
	}
	s.ctx, s.cancel = context.WithCancel(context.TODO())
	err := s.registry.Register(name, net, s.Id, addr)
	if err != nil {
		return nil
	}
	return s
}

// block func
func (r *Service) ListenAndServe(ctx context.Context, router *router.Router) error {
	defer func() {
		r.registry.DeRegister(r.Name, r.Net, r.Id)
	}()
	var err error

	// transport options
	toptions := []transport.Option{}
	if router == nil {
		return fmt.Errorf("router is nil")
	}
	toptions = append(toptions, transport.WithRouter(router))

	if r.Net == "tcp" || r.Net == "tcp4" || r.Net == "tcp6" {
		r.trans, err = transport.NewTcpServerTransport(ctx, r.Net, r.Addr, r.codec, toptions...)
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
