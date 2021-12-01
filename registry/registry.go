package registry

import (
	"sync"

	"github.com/qianxi0410/naive-rpc/server"
)

var (
	registries    = map[string]Registry{}
	registriesLck = sync.RWMutex{}
)

// Registry registry interacts with the remote Nameing Service
type Registry interface {
	// Register, register service
	Register(service *server.Service, opts ...Option) error
	// UnRegister, unregister service
	DeRegister(service *server.Service) error

	// GetService, get services by name, which may have more than one version
	GetService(name string) ([]*server.Service, error)
	// ListServices, list all registered services
	ListServices() ([]*server.Service, error)

	// Watcher, returns a watcher, which watches events on NamingService backend
	Watcher() (Watcher, error)
}

// Option registry option
type Option func(options *options)

type options struct{}

// Watcher watch event from remote Naming Service
type Watcher interface {
	Next() (*Result, error)
	Stop()
}

// Result watch result of event
type Result struct {
	Action ActionType
}

type ActionType = int

const (
	ActionTypeCreate = iota
	ActionTypeUpdate
	ActionTypeDelete
)

func RegisterRegistry(name string, registry Registry) {
	registriesLck.Lock()
	registries[name] = registry
	registriesLck.Unlock()
}

func GetRegistry(name string) Registry {
	registriesLck.RLock()
	defer registriesLck.RUnlock()

	v, ok := registries[name]
	if !ok {
		return nil
	}
	return v
}
