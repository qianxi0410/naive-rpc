package naiverpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithConf(t *testing.T) {
	opts := options{}

	o := WithConf("../loveyourself.ini")
	o(&opts)

	assert.Equal(t, "../loveyourself.ini", opts.conf, "they should be equal")
}
