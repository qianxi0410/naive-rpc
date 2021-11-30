package config

import (
	"fmt"
	"os"
	"sync/atomic"
)

// read config file
type Config interface {
	// if not found return default value
	Read(key, defaultValue string) string

	ReadInt(key string, defaultValue int) int

	ReadBool(key string, defaultValue bool) bool
}

type config struct {
	value atomic.Value
	// opts  options
}

func (r *config) Read(key, defaultValue string) string {
	cfg := r.value.Load()

	switch v := cfg.(type) {
	case *YAMLConfig:
		return v.Read(key, defaultValue)
	default:
		fmt.Fprintf(os.Stderr, "not support config %T\n", v)
		return defaultValue
	}
}

func (r *config) ReadInt(key string, defaultValue int) int {
	cfg := r.value.Load()

	switch v := cfg.(type) {
	case *YAMLConfig:
		return v.ReadInt(key, defaultValue)
	default:
		fmt.Fprintf(os.Stderr, "not support config %T\n", v)
		return defaultValue
	}
}

func (r *config) ReadBool(key string, defaultValue bool) bool {
	cfg := r.value.Load()

	switch v := cfg.(type) {
	case *YAMLConfig:
		return v.ReadBool(key, defaultValue)
	default:
		fmt.Fprintf(os.Stderr, "not support config %T\n", v)
		return defaultValue
	}
}
