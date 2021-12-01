package naiverpc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qianxi0410/naive-rpc/codec/evangelion"
	"github.com/qianxi0410/naive-rpc/config"
	"github.com/qianxi0410/naive-rpc/registry"
	"github.com/qianxi0410/naive-rpc/server"
)

// ListenAndServe quickly initialize Service and ServerModules and serve
func ListenAndServe(opts ...Option) {
	options := options{
		conf: "./conf/service.yaml",
	}
	for _, o := range opts {
		o(&options)
	}

	// load config
	cfg, err := loadConfig(options.conf)
	if err != nil {
		panic(err)
	}

	proc := cfg.Read("name", "eva")
	service := server.NewService(proc)

	tcpport := cfg.ReadInt("tcp_port", 0)

	if err := initTransport(service, "tcp4", tcpport, evangelion.NAME); err != nil {
		panic(err)
	}

	// register to naming service
	registryName := cfg.Read("registry", "eva")
	registry := registry.GetRegistry(registryName)
	if err := registry.Register(service); err != nil {
		panic(err)
	}
}

func loadConfig(fp string) (*config.YAMLConfig, error) {

	if !filepath.IsAbs(fp) {
		self, err := os.Executable()
		if err != nil {
			return nil, err
		}
		dir, _ := filepath.Split(self)
		fp = filepath.Join(dir, fp)
	}

	// load config
	cfg, err := config.NewYamlConfig(fp)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func initTransport(service *server.Service, network string, port int, codec string) error {

	if !(len(network) != 0 &&
		(network == "tcp" || network == "tcp4" || network == "tcp6")) {
		return fmt.Errorf("invalid network: %s", network)
	}

	if port <= 0 {
		return fmt.Errorf("invalid port: %d", port)
	}

	addr := fmt.Sprintf(":%d", port)

	err := service.ListenAndServe(context.Background(), network, addr, codec)
	if err != nil {
		return err
	}
	return nil
}
