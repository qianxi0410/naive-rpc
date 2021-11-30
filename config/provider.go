package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type ProviderType string

const (
	ProviderGFileSystem ProviderType = "fs"
	// maybe etcd consul
)

var providers = map[ProviderType]Provider{
	ProviderGFileSystem: &FileSystemProvider{},
}

type Provider interface {
	Watcher
	Name() string
	Load(ctx context.Context, fp string) ([]byte, error)
}

type FileSystemProvider struct {
}

func (p *FileSystemProvider) Name() string {
	return "fs"
}

func (r *FileSystemProvider) Watch(ctx context.Context, fp string) (<-chan Event, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ch := make(chan Event)
	go func() {
		defer func() {
			watcher.Close()
			close(ch)
		}()

		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					fmt.Println("watch error: wather.Events closed")
					return
				}
				if ev.Op&fsnotify.Write == 0 {
					continue
				}
				ch <- Event{
					typ:  Update,
					meta: ev.String(),
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("watch error:", err)
				return
			case <-ctx.Done():
				log.Println("watcher cancel")
				return
			}
		}
	}()

	err = watcher.Add(fp)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (r *FileSystemProvider) Load(ctx context.Context, fp string) ([]byte, error) {
	fin, err := os.Lstat(fp)
	if err != nil {
		return nil, err
	}

	if fin.IsDir() {
		return nil, errors.New("not normal file")
	}

	return os.ReadFile(fp)
}
