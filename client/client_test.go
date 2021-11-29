package client

import (
	"testing"

	"github.com/qianxi0410/naive-rpc/client/selector"
)

func newClient(name, address, codec string, selector selector.Selector) Client {

	opts := []Option{
		WithAddr(address),
		WithCodec(codec),
		WithSelector(selector),
	}
	client := NewClient(name, opts...)
	return client
}

func TestNewClient(t *testing.T) {
	client := newClient("greeter", "ip://127.0.0.1:8888", "EVANGELION", &selector.IPSelector{})
	t.Logf("client:%+v", client)
}
