package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var fsProvider = &FileSystemProvider{}

func TestFilesystemProviderLoad(t *testing.T) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(d, "testdata/test.yaml")

	b, err := fsProvider.Load(context.TODO(), fp)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	t.Logf("load ok: %s", string(b))
}

func TestFilesystemProviderWatch(t *testing.T) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(d, "testdata/test.yaml")

	ctx, cancel := context.WithCancel(context.TODO())
	ch, err := fsProvider.Watch(ctx, fp)
	if err != nil {
		t.Fatalf("watch error: %v", err)
	}

	go func() {
		defer cancel()

		os.WriteFile(fp, []byte("helloworld0"), 0666)
		time.Sleep(time.Second)
		os.WriteFile(fp, []byte("helloworld1"), 0666)
		time.Sleep(time.Second)
		os.WriteFile(fp, []byte("helloworld2"), 0666)
		time.Sleep(time.Second)
	}()

LOOP:
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				break LOOP
			}
			t.Logf("load ok: %s", ev.meta)
		default:
		}
	}
}
