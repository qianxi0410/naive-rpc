package transport

import (
	"context"
	ers "errors"
	"fmt"
	"net"
	"time"

	"github.com/qianxi0410/naive-rpc/client/pool"
	"github.com/qianxi0410/naive-rpc/codec"
	"github.com/qianxi0410/naive-rpc/errors"
)

// send tcp req
type TcpTransport struct {
	Pool  pool.PoolFactory
	Codec codec.Codec
}

// send reqHead and return rspHead, return an error if encountered
func (r *TcpTransport) Send(ctx context.Context, network, addr string, reqHead interface{}) (rsp interface{}, err error) {
	data, err := r.Codec.Encode(reqHead)
	if err != nil {
		return nil, err
	}

	// conn, err := r.Pool.Get(ctx, network, addr)
	// TODO: goroutine pool
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Millisecond * 200))

	n, err := conn.Write(data)
	if err != nil {
		return nil, err
	}

	if len(data) != n {
		return nil, fmt.Errorf("write error, write only %d bytes, want write %d bytes", n, data)
	}

	// allocate buffer
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	// conn read
	var size int
	for {
		n, err := conn.Read(buf[size:])
		if err != nil {
			return nil, err
		}

		size += n
		// decode
		rsp, _, err := r.Codec.Decode(buf[:size])
		if err != nil {
			if ers.Is(err, errors.CodecReadIncompleteErr) {
				continue
			}
			return nil, err
		}

		return rsp, nil
	}
}
