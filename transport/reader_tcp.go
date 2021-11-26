package transport

import (
	"net"
	"sync"
	"time"

	"github.com/qianxi0410/naive-rpc/codec"
	"github.com/qianxi0410/naive-rpc/errors"
)

var tcpBufferPool = &sync.Pool{
	New: func() interface{} {
		// 128kB
		return make([]byte, 128*1024)
	},
}

// TcpMessageReader read req from `net.conn`, if read successfully, return the req'svr session.
//
// if any error occurs, it returns nil session and error, error should be one of the following:
// - io.Timeout
// - ...
type TcpMessageReader struct {
	codec codec.Codec
}

func NewTcpMessageReader(codec codec.Codec) *TcpMessageReader {
	return &TcpMessageReader{codec: codec}
}

func (r *TcpMessageReader) Read(ep *TcpEndPoint) error {
	defer func() {
		ep.conn.Close()
		tcpBufferPool.Put(ep.buf)
		close(ep.reqCh)
	}()

	var (
		bufLen   int
		readSize int
		err      error
	)

	for {
		select {
		case <-ep.ctx.Done():
			return errors.ServerCtxDoneErr
		default:
			// do nothing
		}

		ep.conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		if readSize, err = ep.conn.Read(ep.buf[bufLen:]); err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				time.Sleep(time.Microsecond * 10)
				continue
			}
			return err
		}

		bufLen += readSize

		req, size, err := r.codec.Decode(ep.buf[:bufLen])
		if err != nil {
			if err == errors.CodecReadIncompleteErr {
				continue
			}
			return err
		}
		ep.reqCh <- req
		ep.buf = ep.buf[size:]
		bufLen -= size
	}
}
