package transport

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/qianxi0410/naive-rpc/codec"
)

// tcp endpoint struct
type TcpEndPoint struct {
	conn  net.Conn
	reqCh chan interface{}
	rspCh chan interface{}

	reader *TcpMessageReader
	ctx    context.Context
	cancel context.CancelFunc

	buf []byte
}

func (r *TcpEndPoint) Read() {
	defer func() {
		r.conn.Close()
		r.cancel()
	}()

	err := r.reader.Read(r)
	if err != nil {
		if err != io.EOF {
			log.Printf("peer connection Closed now, local:%s->remote:%s", r.conn.LocalAddr(), r.conn.RemoteAddr())
			return
		}
		log.Fatalf("tcp read request error:%v", err)
	}
}

func (r *TcpEndPoint) Write() {
	defer func() {
		r.conn.Close()
	}()

	for {
		select {
		case <-r.ctx.Done():
			return
		case v := <-r.rspCh:
			session := v.(codec.Session)
			rsp := session.Response()
			data, err := r.reader.codec.Encode(rsp)
			if err != nil {
				log.Fatalf("tcp encode respone error:%v", err)
				continue
			}

			r.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 2500))
			n, err := r.conn.Write(data)
			if err != nil || len(data) != n {
				// fixme handle error
				log.Fatalf("tcp send response error:%v, bytes written got:%d, want:%d", err, n, len(data))
				continue
			}
		}
	}
}
