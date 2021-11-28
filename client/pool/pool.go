package pool

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

type PoolFactory interface {
	Get(ctx context.Context, network, addr string) (net.Conn, error)
}

// create a new poolfactory manager
func NewConnPoolFactory(opt ...Option) PoolFactory {
	opts := &Options{
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		DialTimeout: 200 * time.Millisecond,
	}

	for _, o := range opt {
		o(opts)
	}

	return &poolFactory{
		opts:      opts,
		connPools: &sync.Map{},
	}
}

// poolFactory poolFactory manager, it maintains many <address,Pool> pairs
type poolFactory struct {
	opts      *Options
	connPools *sync.Map
}

func (r *poolFactory) Get(ctx context.Context, network, addr string) (net.Conn, error) {
	var cancel context.CancelFunc

	_, ok := ctx.Deadline()
	if !ok {
		ctx, cancel = context.WithTimeout(ctx, r.opts.DialTimeout)
		defer cancel()
	}
	key := addr + "/" + network

	if v, ok := r.connPools.Load(key); ok {
		return v.(*ConnPool).Get(ctx)
	}

	newPool := &ConnPool{
		dialFunc: func(c context.Context) (net.Conn, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			timeout := r.opts.DialTimeout
			t, ok := ctx.Deadline()
			if ok {
				timeout = t.Sub(time.Now())
			}
			return net.DialTimeout(network, addr, timeout)
		},
		MinIdle:         r.opts.MinIdle,
		MaxIdle:         r.opts.MaxIdle,
		MaxActive:       r.opts.MaxActive,
		Wait:            r.opts.Wait,
		MaxConnLifetime: r.opts.MaxConnLifetime,
		IdleTimeout:     r.opts.IdleTimeout,
		idle:            doublylinkedlist.New(),
	}

	v, ok := r.connPools.LoadOrStore(key, newPool)
	if !ok {
		go newPool.Prepare(ctx)
		newPool.RegisterCheckFunc(time.Second*3, newPool.CheckAlive)
		return newPool.Get(ctx)
	}

	return v.(*ConnPool).Get(ctx)
}
