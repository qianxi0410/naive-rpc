package naiverpc_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	naiverpc "github.com/qianxi0410/naive-rpc"
	_ "github.com/qianxi0410/naive-rpc/codec/evangelion"
)

func TestListenAndServe(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fp := filepath.Join(dir, "conf/service.yaml")

	naiverpc.ListenAndServe(naiverpc.WithConf(fp))

	// Linux: run `fuser port/tcp` or `fuser port/udp` to check whether server working
	// macOS: run `lsof -i tcp:port` or `lsof -i udp:port` to check whether server working
	time.Sleep(time.Second * 10)
}
