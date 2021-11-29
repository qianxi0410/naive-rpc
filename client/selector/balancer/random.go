package balancer

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type RandomBanlancer struct {
	Addrs []string
}

func (r *RandomBanlancer) Next() string {
	idx := rand.Int() % len(r.Addrs)
	return r.Addrs[idx]
}
