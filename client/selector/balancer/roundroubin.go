package balancer

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"go.etcd.io/etcd/clientv3"
)

type RoundRobinBalancer struct {
	Addrs []string
	Cli   *clientv3.Client
}

func (r *RoundRobinBalancer) Next() string {
	var rIdx int64
	rsp, _ := r.Cli.Get(context.TODO(), "RAOUND_ROUBIN")

	if len(rsp.Kvs) != 0 {
		i, _ := strconv.Atoi(string(rsp.Kvs[0].Value))
		rIdx = int64(i)
	} else {
		rIdx = -1
	}

	idx := atomic.AddInt64(&rIdx, 1)
	idx = idx % int64(len(r.Addrs))
	r.Cli.Put(context.TODO(), "RAOUND_ROUBIN", fmt.Sprintf("%d", rIdx))
	return r.Addrs[idx]
}
