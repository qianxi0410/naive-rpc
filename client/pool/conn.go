package pool

import (
	"context"
	ers "errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/qianxi0410/naive-rpc/errors"
)

type ConnPool struct {
	dialFunc  func(context.Context) (net.Conn, error)
	MinIdle   int
	MaxIdle   int
	MaxActive int
	Wait      bool
	// ch is init
	initialized uint32
	// when wait is true, to limit the conn number
	ch              chan struct{}
	IdleTimeout     time.Duration
	MaxConnLifetime time.Duration
	// lock
	connMutex sync.Mutex
	// represent pool cloesd ?
	closed bool
	// current active conn
	active int
	// idle conn list
	idle *doublylinkedlist.List
}

// get a conn from pool
func (r *ConnPool) Get(ctx context.Context) (ci *ConnItem, err error) {
	for {
		if ci, err = r.get(ctx); err != nil {
			return nil, err
		}

		if ci.readClosed() {
			r.put(ci, true)
			continue
		}
		return ci, nil
	}
}

// close the pool
func (r *ConnPool) Close() error {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true
	r.active = r.idle.Size()

	if r.ch != nil {
		close(r.ch)
	}

	for i := 0; i < r.idle.Size(); i++ {
		item, ok := r.idle.Get(i)
		if !ok {
			break
		}
		ci := item.(*ConnItem)
		ci.Conn.Close()
		ci.closed = true
	}

	return nil
}

// Prepare prepare MinIdle number of connections in advance
// so can reduce the propability
// of creating connections when launch IO actions.
func (r *ConnPool) Prepare(ctx context.Context) {
	if r.MinIdle <= 0 {
		return
	}

	if r.MinIdle > r.MaxIdle {
		r.MinIdle = r.MaxIdle
	}

	if r.MinIdle > r.MaxActive {
		r.MinIdle = r.MaxActive
	}

	if r.Wait && r.MaxActive > 0 {
		r.initializeCh()
	}

	conns := make([]*ConnItem, 0, r.MinIdle)
	for i := 0; i < r.MinIdle; i++ {
		for {
			poolConn, err := r.get(ctx)
			if err != nil {
				continue
			}
			conns = append(conns, poolConn)
			break
		}
	}

	for _, poolConn := range conns {
		r.put(poolConn, poolConn.readClosed())
	}
}

// instance in conn pool
type ConnItem struct {
	net.Conn
	recycled time.Time
	created  time.Time
	pool     *ConnPool
	closed   bool
}

func (r *ConnItem) reset() {
	if r == nil {
		return
	}
	r.Conn.SetDeadline(time.Time{})
}

func (r *ConnItem) Write(data []byte) (int, error) {
	if r.closed {
		return 0, errors.ConnClosedErr
	}
	n, err := r.Conn.Write(data)
	if err != nil {
		r.pool.put(r, true)
	}

	return n, nil
}

func (r *ConnPool) initializeCh() {
	if atomic.LoadUint32(&r.initialized) == 1 {
		return
	}
	r.connMutex.Lock()

	if r.initialized == 0 {
		r.ch = make(chan struct{}, r.MaxActive)
		if r.closed {
			close(r.ch)
		} else {
			for i := 0; i < r.MaxActive; i++ {
				r.ch <- struct{}{}
			}
		}
		atomic.StoreUint32(&r.initialized, 1)
	}
	r.connMutex.Unlock()
}

// exceedLimit check whether number of connections has reached the limit
func (r *ConnPool) exceedLimit() bool {
	return !r.Wait && r.MaxActive > 0 && r.active >= r.MaxActive
}

// dial create a new connection
func (r *ConnPool) dial(ctx context.Context) (net.Conn, error) {
	if r.dialFunc != nil {
		return r.dialFunc(ctx)
	}
	return nil, ers.New("must pass dialFunc to poolFactory")
}

func (r *ConnPool) get(ctx context.Context) (*ConnItem, error) {
	// if `wait`, then should initialize a `chan` to sync
	if r.Wait && r.MaxActive > 0 {
		r.initializeCh()

		select {
		case <-r.ch:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	if r.closed {
		return nil, errors.PoolClosedErr
	}

	if r.exceedLimit() {
		return nil, errors.ExceedPoolLimitErr
	}

	v, ok := r.idle.Get(0)

	if ok && v != nil {
		r.idle.Remove(0)
		ci := v.(*ConnItem)
		return ci, nil
	}

	conn, err := r.dial(ctx)
	if err != nil {
		if r.ch != nil && !r.closed {
			r.ch <- struct{}{}
		}
		return nil, err
	}
	r.active++

	return &ConnItem{Conn: conn, created: time.Now(), pool: r}, nil
}

// RegisterCheckFunc register function to check whether a connection is alive
func (r *ConnPool) RegisterCheckFunc(interval time.Duration, checkFunc func(*ConnItem) bool) {
	if interval <= 0 || checkFunc == nil {
		return
	}

	go func() {
		for {
			time.Sleep(interval)
			r.connMutex.Lock()
			size := r.idle.Size()
			r.connMutex.Unlock()

			var ci *ConnItem
			for i := 0; i < size; i++ {
				r.connMutex.Lock()
				v, ok := r.idle.Get(0)
				if ok && v != nil {
					r.idle.Remove(0)
					ci = v.(*ConnItem)
				}

				if !checkFunc(ci) {
					ci.Conn.Close()
					ci.closed = true
					r.connMutex.Lock()
					r.active--
					r.connMutex.Unlock()
				} else {
					r.connMutex.Lock()
					r.idle.Add(ci)
					r.connMutex.Unlock()
				}
			}
		}
	}()
}

// CheckAlive default checkfunc to test whether a connection is alive or not
func (r *ConnPool) CheckAlive(ci *ConnItem) bool {
	// check whether connection is idle
	if r.IdleTimeout > 0 && ci.recycled.Add(r.IdleTimeout).Before(time.Now()) {
		return true
	}
	// check whether connection lifecyle is ok
	if r.MaxConnLifetime > 0 && ci.created.Add(r.MaxConnLifetime).Before(time.Now()) {
		return true
	}

	// check whether read half closed
	return ci.readClosed()
}

// put try put the connection back into poolFactory
func (r *ConnPool) put(ci *ConnItem, forceClose bool) error {
	r.connMutex.Lock()
	defer r.connMutex.Unlock()

	if !r.closed && !forceClose {
		ci.recycled = time.Now()
		r.idle.Insert(0, ci)
		if r.idle.Size() > r.MaxIdle {
			var item *ConnItem
			lastIdle, ok := r.idle.Get(r.idle.Size() - 1)
			if ok {
				r.idle.Remove(r.idle.Size() - 1)
				item = lastIdle.(*ConnItem)
				item.closed = true
				item.Conn.Close()
				r.active--
			}
		}
	}

	if r.Wait && r.ch != nil && !r.closed {
		r.ch <- struct{}{}
	}

	return nil
}

func (r *ConnItem) Read(buf []byte) (int, error) {
	if r.closed {
		return 0, errors.ConnClosedErr
	}
	n, err := r.Conn.Read(buf)
	if err != nil {
		r.pool.put(r, true)
	}
	return n, nil
}

func (r *ConnItem) Close() error {
	if r.closed {
		return errors.ConnClosedErr
	}
	r.reset()
	return r.pool.put(r, false)
}

func (r *ConnItem) readClosed() bool {
	return readClosed(r.Conn)
}
