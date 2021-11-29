package selector

import "github.com/qianxi0410/naive-rpc/client/selector/balancer"

type IPSelector struct {
	network  string
	addrs    []string
	balancer balancer.Balancer
}

func NewIPSelector(network string, addrs []string) *IPSelector {
	return &IPSelector{
		network:  network,
		addrs:    addrs,
		balancer: &balancer.RandomBanlancer{Addrs: addrs},
	}
}

func (r *IPSelector) Select(service string) (*Node, error) {
	addr := r.balancer.Next()
	node := Node{
		Network: r.network,
		Address: addr,
	}
	return &node, nil
}

// Update do nothing in IPSelector
func (s *IPSelector) Update(node *Node, err error) error {
	return nil
}
