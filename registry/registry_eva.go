package registry

import (
	"fmt"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"go.etcd.io/etcd/clientv3"
)

func init() {
	RegisterRegistry("eva", &EvaRegistry{})
}

type EvaRegistry struct {
	cli *clientv3.Client
}

func NewEvaRegistry(c clientv3.Config) *EvaRegistry {
	cli, err := clientv3.New(c)
	if err != nil {
		return nil
	}

	return &EvaRegistry{
		cli: cli,
	}
}

func (r *EvaRegistry) Register(name, net, id, addr string, opts ...Option) error {
	key := fmt.Sprintf("%s/%s/%s", name, net, id)
	_, err := r.cli.Put(context.TODO(), key, addr)
	return err
}

func (r *EvaRegistry) DeRegister(name, net, id string) error {
	key := fmt.Sprintf("%s/%s/%s", name, net, id)
	_, err := r.cli.Delete(context.TODO(), key)
	return err
}

func (r *EvaRegistry) GetAddrs(name, net string) ([]string, error) {

	keyPrefix := fmt.Sprintf("%s/%s/", name, net)
	kvs, err := r.cli.Get(context.TODO(), keyPrefix, clientv3.WithPrefix())

	if err != nil {
		return nil, err
	}
	addrs := make([]string, 0, len(kvs.Kvs))

	for _, kv := range kvs.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	return addrs, nil
}

func (r *EvaRegistry) Watcher() (Watcher, error) {
	return nil, nil
}
