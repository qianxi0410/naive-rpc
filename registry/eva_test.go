package registry_test

import (
	"log"
	"testing"
	"time"

	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"github.com/qianxi0410/naive-rpc/registry"
	"github.com/qianxi0410/naive-rpc/server"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"
)

var eva registry.EvaRegistry

func init() {
	eva = *registry.NewEvaRegistry(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:3000"},
		DialTimeout: 5 * time.Second,
	})
}

func TestRegister(t *testing.T) {
	s1 := server.NewService("test2", "tcp", "127.0.0.1:8888", evangelion.NAME)
	s2 := server.NewService("test2", "tcp", "127.0.0.1:9999", evangelion.NAME)

	err := eva.Register(s1)
	assert.Nil(t, err)
	err = eva.Register(s2)
	assert.Nil(t, err)
}

func TestGetAddrs(t *testing.T) {
	s, err := eva.GetAddrs("test2", "tcp")
	assert.Nil(t, err)
	assert.Equal(t, len(s), 2)
	log.Println(s)
}

func TestDeRegistry(t *testing.T) {
	s1 := server.NewService("test2", "tcp", "127.0.0.1:8888", evangelion.NAME)
	eva.DeRegister(s1)
}
