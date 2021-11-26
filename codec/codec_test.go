package codec_test

import (
	"testing"

	"github.com/qianxi0410/naive-rpc/codec"
	"github.com/stretchr/testify/assert"
)

type fakeCodec struct{}

func (r *fakeCodec) Name() string {
	return "fake"
}

func (r *fakeCodec) Encode(pkg interface{}) (data []byte, err error) {
	return nil, nil
}

func (r *fakeCodec) Decode(data []byte) (req interface{}, n int, err error) {
	return nil, -1, nil
}

func TestRegisterCodec(t *testing.T) {
	f := &fakeCodec{}
	codec.RegisterCodec("fake", f, f)
	assert.Equal(t, codec.ServerCodec("fake"), f)
	assert.Equal(t, codec.ClientCodec("fake"), f)
}
