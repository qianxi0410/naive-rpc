package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type yamlStruct struct {
	Name  string  `yaml:"name"`
	Age   int     `yaml:"age"`
	Float float64 `yaml:"float"`
	Recv  struct {
		B struct {
			C int `yaml:"c"`
		} `yaml:"b"`
	} `yaml:"recv"`
}

func TestYAMLDecoderDecode(t *testing.T) {
	d := &YAMLDecoder{}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := filepath.Join(cwd, "testdata/test.yaml")
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}

	// v := &yamlStruct{}
	cfg, err := d.Decode(b)
	if err != nil {
		t.Fatalf("yaml decode error: %v", err)
	}
	t.Logf("yaml decode ok, data: %s", cfg.(*YAMLConfig).Read("name", ""))
}
