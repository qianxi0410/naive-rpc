package codec

import (
	"bytes"
	"compress/gzip"
	"io"
)

type Compressor interface {
	// compress data
	Compress(data []byte) ([]byte, error)
	// decomress data
	Decompress(data []byte) ([]byte, error)
}

// compress type
type CompressType int

const (
	CompressGZip = CompressType(iota)
)

type GZipCompressor struct{}

func (r *GZipCompressor) Compress(data []byte) ([]byte, error) {
	buf := &bytes.Buffer{}
	w, err := gzip.NewWriterLevel(buf, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

func (r *GZipCompressor) Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	return io.ReadAll(reader)
}
