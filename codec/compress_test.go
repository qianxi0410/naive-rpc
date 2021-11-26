package codec_test

import (
	"encoding/hex"
	"testing"

	"github.com/qianxi0410/naive-rpc/codec"
	"github.com/stretchr/testify/assert"
)

var (
	msg = "helloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworld"

	compressor = &codec.GZipCompressor{}
)

func TestGZipCompressorCompress(t *testing.T) {
	b, err := compressor.Compress([]byte(msg))
	assert.Nil(t, err)
	t.Logf("compressed data hex: %s, origin len %d, len: %d", hex.EncodeToString(b), len([]byte(msg)), len(b))
}

func TestGZipCompressorDecompress(t *testing.T) {
	b, err := compressor.Compress([]byte(msg))
	assert.Nil(t, err)
	assert.Len(t, b, 55)

	data, err := compressor.Decompress(b)
	assert.Nil(t, err)
	c := string(data)
	assert.Equal(t, msg, c)
}
