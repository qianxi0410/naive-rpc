package registry

import (
	"log"
	"testing"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"
)

func TestPut(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:3000"},
		DialTimeout: time.Second * 5,
	})

	assert.Nil(t, err)
	_, err = cli.Put(context.TODO(), "key", "value")
	assert.Nil(t, err)

}

func TestPutDir(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:3000"},
		DialTimeout: time.Second * 5,
	})
	kv := clientv3.NewKV(cli)
	assert.Nil(t, err)
	_, err = kv.Put(context.TODO(), "/keys/key1", "value1")
	assert.Nil(t, err)
	_, err = kv.Put(context.TODO(), "/keys/key2", "value2")
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:3000"},
		DialTimeout: time.Second * 5,
	})
	kv := clientv3.NewKV(cli)
	assert.Nil(t, err)
	rsp, err := kv.Get(context.TODO(), "/keys/", clientv3.WithPrefix())
	assert.Nil(t, err)
	for _, kv := range rsp.Kvs {
		log.Printf("%s: %s", kv.Key, kv.Value)
	}
}

func TestWatch(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:3000"},
		DialTimeout: time.Second * 5,
	})
	assert.Nil(t, err)

	wc := cli.Watch(context.TODO(), "/keys", clientv3.WithPrefix())
	for c := range wc {
		for _, ev := range c.Events {
			log.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func TestDelete(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:3000"},
		DialTimeout: time.Second * 5,
	})
	assert.Nil(t, err)

	cli.Delete(context.TODO(), "key")
}
