package selector

import (
	"github.com/kpango/glg"
	"github.com/qianxi0410/naive-rpc/client/selector/balancer"
	"github.com/qianxi0410/naive-rpc/registry"
	"go.etcd.io/etcd/clientv3"
)

type SelectorType int

const (
	RANDOM SelectorType = iota
	ROUND_ROUBIN
)

type IPSelector struct {
	name     string
	network  string
	addrs    []string
	registry registry.Registry
	balancer balancer.Balancer
}

func NewIPSelector(name, network string, typ SelectorType, conf clientv3.Config) Selector {
	var selector IPSelector

	selector.name = name
	selector.network = network
	selector.registry = registry.NewEvaRegistry(conf)

	addrs, err := selector.registry.GetAddrs(name, network)
	if err != nil {
		return nil
	}
	selector.addrs = append(selector.addrs, addrs...)

	switch typ {
	case RANDOM:
		selector.balancer = &balancer.RandomBanlancer{Addrs: selector.addrs}
		glg.Success("RANDOM strategy selector is init")
	case ROUND_ROUBIN:
		cli, _ := clientv3.New(conf)
		selector.balancer = &balancer.RoundRobinBalancer{Addrs: selector.addrs, Cli: cli}
		glg.Success("ROUND_ROUBIN strategy selector is init")
	default:
		selector.balancer = &balancer.RandomBanlancer{Addrs: selector.addrs}
		glg.Success("RANDOM strategy selector is init")
	}

	return &selector
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
