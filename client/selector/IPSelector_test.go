package selector_test

import (
	"testing"

	"github.com/qianxi0410/naive-rpc/client/selector"
	"go.etcd.io/etcd/clientv3"
)

func TestIPSelector(t *testing.T) {
	var (
		network = "tcp"
		address = "127.0.0.1:8888"
	)
	s := selector.NewIPSelector("test2", "tcp", selector.RANDOM, clientv3.Config{
		Endpoints: []string{"127.0.0.1:3000"},
	})
	if s == nil {
		t.Fatalf("ipselector create failed")
	}

	n, err := s.Select("")
	if err != nil {
		t.Fatalf("ipselector select error:%v", err)
	}

	if n.Network != "tcp" && n.Address != "127.0.0.1:9999" {
		t.Fatalf("ipselector select error, got %s/%s, want %s/%s", n.Address, n.Network, address, network)
	}
}
