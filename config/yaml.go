package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/smallfish/simpleyaml"
)

type YAMLConfig struct {
	yml *simpleyaml.Yaml
}

func NewYamlConfig(fp string) (*YAMLConfig, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	yml, err := simpleyaml.NewYaml(data)
	if err != nil {
		return nil, err
	}

	return &YAMLConfig{yml: yml}, nil
}

func (r *YAMLConfig) path(key string) []interface{} {
	path := strings.Split(key, ".")
	paths := make([]interface{}, len(path))
	for idx, p := range path {
		paths[idx] = p
	}
	return paths
}

func (r *YAMLConfig) Read(key, defaultValue string) string {
	if r.yml == nil {
		return defaultValue
	}

	path := r.path(key)
	v, err := r.yml.GetPath(path...).String()
	if err != nil || len(v) == 0 {
		return defaultValue
	}

	return v
}

func (r *YAMLConfig) ReadInt(key string, defaultValue int) int {
	if r.yml == nil {
		return defaultValue
	}

	path := r.path(key)
	v, err := r.yml.GetPath(path...).Int()
	if err != nil || v == 0 {
		return defaultValue
	}

	return v
}

func (r *YAMLConfig) ReadBool(key string, defaultValue bool) bool {
	if r.yml == nil {
		return defaultValue
	}

	path := r.path(key)
	v, err := r.yml.GetPath(path...).Bool()
	if err != nil {
		return defaultValue
	}

	return v
}
