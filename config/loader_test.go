package config_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/qianxi0410/naive-rpc/config"
	"github.com/stretchr/testify/assert"
)

func TestLoaderLoad(t *testing.T) {
	opts := []config.Option{
		config.WithProvider(&config.FileSystemProvider{}),
		config.WithDecoder(&config.YAMLDecoder{}),
	}
	ld, err := config.NewLoader(context.TODO(), opts...)
	if err != nil {
		t.Fatalf("new loader error: %v", err)
	}

	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c, err := ld.Load(context.TODO(), filepath.Join(d, "testdata/test.yaml"))
	if err != nil {
		t.Fatalf("load test.yml error: %v", err)
	}

	assert.Equal(t, "qianxi", c.Read("name", ""))
	assert.Equal(t, 10, c.ReadInt("recv/a/b", 10))
}

func TestLoaderReLoad(t *testing.T) {
	opts := []config.Option{
		config.WithProvider(&config.FileSystemProvider{}),
		config.WithDecoder(&config.YAMLDecoder{}),
		config.WithReload(true),
	}
	ld, err := config.NewLoader(context.TODO(), opts...)
	if err != nil {
		t.Fatalf("new loader error: %v", err)
	}

	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(d, "testdata/test.yaml")
	c, err := ld.Load(context.TODO(), fp)
	if err != nil {
		t.Fatalf("load test.yaml error: %v", err)
	}

	assert.Equal(t, "qianxi", c.Read("name", ""))

	// change file
	b, err := os.ReadFile(fp)
	if err != nil {
		panic(err)
	}
	defer os.WriteFile(fp, b, 0666)

	n := strings.ReplaceAll(string(b), "name: qianxi", "name: qianxi2")
	err = os.WriteFile(fp, []byte(n), 0666)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 3)
	assert.Equal(t, "qianxi2", c.Read("name", ""))
}
