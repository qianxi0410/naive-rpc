package registry

import "github.com/qianxi0410/naive-rpc/server"

func init() {
	RegisterRegistry("eva", &EvaRegistry{})
}

type EvaRegistry struct{}

func (r *EvaRegistry) Register(service *server.Service, opts ...Option) error {
	return nil
}

func (r *EvaRegistry) DeRegister(service *server.Service) error {
	return nil
}

func (r *EvaRegistry) GetService(name string) ([]*server.Service, error) {
	return nil, nil
}

func (r *EvaRegistry) ListServices() ([]*server.Service, error) {
	return nil, nil
}

func (r *EvaRegistry) Watcher() (Watcher, error) {
	return nil, nil
}
