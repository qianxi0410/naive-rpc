package config

import (
	"github.com/smallfish/simpleyaml"
)

type DecoderType int

const (
	YAML = DecoderType(iota)
)

type Decoder interface {
	// decode data into val, val should be pointer
	Decode(data []byte) (Config, error)
}

type YAMLDecoder struct{}

func (r *YAMLDecoder) Decode(data []byte) (Config, error) {
	// rt := reflect.TypeOf(val)
	// if rt.Kind() != reflect.Ptr && rt.Kind() != reflect.Map {
	// 	panic("val must be a pointer")
	// }
	// err := yaml.Unmarshal(data, val)
	// return err
	y, err := simpleyaml.NewYaml(data)
	if err != nil {
		return nil, err
	}

	return &YAMLConfig{
		yml: y,
	}, nil
}
