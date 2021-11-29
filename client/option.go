package client

import (
	"github.com/qianxi0410/naive-rpc/client/selector"
	"github.com/qianxi0410/naive-rpc/client/transport"
	"github.com/qianxi0410/naive-rpc/codec"
)

type Option func(*client)

func WithAddr(addr string) Option {
	return func(r *client) {
		r.Addr = addr
	}
}

// TransportType options
type TransportType int

const (
	UDP = iota
	UDP4
	UDP6
	TCP
	TCP4
	TCP6
	UNIX
)

func (t TransportType) String() string {
	switch t {
	case UDP:
		return "udp"
	case UDP4:
		return "udp4"
	case UDP6:
		return "udp6"
	case TCP:
		return "tcp"
	case TCP4:
		return "tcp4"
	case TCP6:
		return "tcp6"
	case UNIX:
		return "unix"
	default:
		return ""
	}
}

func (t TransportType) Valid() bool {
	if t == UDP || t == UDP4 || t == UDP6 ||
		t == TCP || t == TCP4 || t == TCP6 ||
		t == UNIX {
		return true
	}
	return false
}

// WithTransportType specify the transport type, support UDP, TCP, Unix
func WithTransportType(typ TransportType) Option {
	return func(r *client) {
		r.TransportType = typ
		switch typ {
		case TCP, TCP4, TCP6:
			r.Transport = &transport.TcpTransport{
				Pool:  defaultPoolFactory,
				Codec: codec.ClientCodec("whisper"),
			}
		}
	}
}

// RpcType options
type RpcType int

const (
	SendOnly = iota
	SendRecv
	SendRecvMultiplex
	SendStreamOnly
	SendStreamAndRecv
	SendAndRecvStream
	SendStreamAndRecvStream
)

// WithRpcType specify the rpc type, support SendOnly, SendRecv, SendRecvWithMultiplex, etc.
func WithRpcType(typ RpcType) Option {
	return func(r *client) {
		r.RpcType = typ
	}
}

// WithSelector specify the selector
func WithSelector(selector selector.Selector) Option {
	return func(r *client) {
		r.Selector = selector
	}
}

// WithCodec specify the codec
func WithCodec(name string) Option {
	return func(r *client) {
		r.Codec = codec.ClientCodec(name)
	}
}
