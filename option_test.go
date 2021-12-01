package naiverpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithConf(t *testing.T) {
	opts := options{}

	o := WithConf("../loveyourself.yaml")
	o(&opts)

	assert.Equal(t, "../loveyourself.yaml", opts.conf, "they should be equal")
}
