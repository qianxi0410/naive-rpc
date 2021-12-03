package server

import (
	"context"
	"sync"

	"github.com/qianxi0410/naive-rpc/router"
	"github.com/qianxi0410/naive-rpc/transport"
)

// Service represents a server instance (a server process),
// Any server can include more than one service, i.e, any server can be
// plugged into multile modules, like TcpServerTransport, UdpServerTransport, Broker, etc.
// By this way, we can implement more modules to extend server's abilities.
type Service struct {
	name   string
	ctx    context.Context
	cancel context.CancelFunc
	opts   *options

	trans      []transport.Transport
	transMutex *sync.Mutex

	router    *router.Router
	startOnce sync.Once
	stopOnce  sync.Once
	closed    chan (struct{})
}

// create a new server with option
func NewService(name string, opts ...Option) *Service {
	s := &Service{
		name:       name,
		opts:       &options{},
		trans:      []transport.Transport{},
		transMutex: &sync.Mutex{},
		router:     router.NewRouter(),
		startOnce:  sync.Once{},
		stopOnce:   sync.Once{},
		closed:     make(chan struct{}, 1),
	}
	s.ctx, s.cancel = context.WithCancel(context.TODO())

	for _, o := range opts {
		o(s.opts)
	}
	return s
}

func (s *Service) ListenAndServe(ctx context.Context, net, addr, codec string, opts ...Option) error {
	var (
		trans transport.Transport
		err   error
	)

	options := options{}
	for _, o := range opts {
		o(&options)
	}

	// transport options
	toptions := []transport.Option{}
	if options.Router != nil {
		toptions = append(toptions, transport.WithRouter(options.Router))
	}

	if net == "tcp" || net == "tcp4" || net == "tcp6" {
		trans, err = transport.NewTcpServerTransport(ctx, net, addr, codec, toptions...)
		if err != nil {
			return err
		}

	}

	s.transMutex.Lock()
	s.trans = append(s.trans, trans)
	s.transMutex.Unlock()
	go trans.ListenAndServe()

	select {
	case <-ctx.Done():
		return nil
	case <-trans.Closed():
		return nil
	}

}
