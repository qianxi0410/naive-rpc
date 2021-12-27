package evangelion

import (
	"bytes"
	"encoding/binary"
	ers "errors"
	"fmt"

	"github.com/qianxi0410/naive-rpc/errors"
	"google.golang.org/protobuf/proto"
)

const maxDefaultPkgSize = 10 * (1 << 20) // 10 MB

// server side codec
type ServerCodec struct{}

func (s *ServerCodec) Name() string {
	return NAME
}

func (s *ServerCodec) Encode(pkg interface{}) ([]byte, error) {

	pb, ok := pkg.(*Response)
	if !ok {
		return nil, ers.New("pkg not valid *internal.Response")
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	binary.Write(b, binary.BigEndian, int8(0x38))
	binary.Write(b, binary.BigEndian, int32(len(data)))
	binary.Write(b, binary.BigEndian, data)
	binary.Write(b, binary.BigEndian, int8(0x49))

	return b.Bytes(), nil
}

func (s *ServerCodec) Decode(in []byte) (interface{}, int, error) {

	if len(in) < 5 {
		return nil, 0, errors.CodecReadIncompleteErr
	}

	b := bytes.NewBuffer(in)

	// pkg: | 1B:0x38 | 4B:len | payload | 1B: 0x49 |
	var (
		pkgStx int8
		pkgLen int32
		pkgEtx int8
	)

	// stx
	if err := binary.Read(b, binary.BigEndian, &pkgStx); err != nil {
		return nil, 0, fmt.Errorf("read stx: %v", err)
	}
	if pkgStx != 0x38 {
		return nil, 0, fmt.Errorf("%w: stx mismatch", errors.CodecReadInvalid)
	}

	// len
	if err := binary.Read(b, binary.BigEndian, &pkgLen); err != nil {
		return nil, 0, fmt.Errorf("read len: %v", err)
	}
	if pkgLen > maxDefaultPkgSize {
		return nil, 0, errors.CodecRead2BigErr
	}

	totalLen := int(1 + 4 + pkgLen + 1)
	if len(in) < totalLen {
		return nil, 0, errors.CodecReadIncompleteErr
	}

	// payload
	payload := make([]byte, pkgLen)
	if err := binary.Read(b, binary.BigEndian, payload); err != nil {
		return nil, 0, errors.CodecReadErr
	}

	// etx
	if err := binary.Read(b, binary.BigEndian, &pkgEtx); err != nil {
		return nil, 0, fmt.Errorf("read ext: %v", err)
	}
	if pkgEtx != 0x49 {
		return nil, 0, fmt.Errorf("%w: etx mismatch", errors.CodecReadInvalid)
	}

	request := &Request{}
	if err := proto.Unmarshal(payload, request); err != nil {
		return nil, 0, err
	}

	return request, totalLen, nil
}

// ClientCodec clientside codec
type ClientCodec struct{}

func (c *ClientCodec) Name() string {
	return NAME
}

func (c *ClientCodec) Encode(pkg interface{}) ([]byte, error) {

	pb, ok := pkg.(*Request)
	if !ok {
		return nil, ers.New("pkg not valid *evangelion.ReqHead")
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}

	b := &bytes.Buffer{}
	binary.Write(b, binary.BigEndian, int8(0x38))
	binary.Write(b, binary.BigEndian, int32(len(data)))
	binary.Write(b, binary.BigEndian, data)
	binary.Write(b, binary.BigEndian, int8(0x49))

	return b.Bytes(), nil
}

func (c *ClientCodec) Decode(in []byte) (interface{}, int, error) {

	if len(in) < 5 {
		return nil, 0, errors.CodecReadIncompleteErr
	}

	b := bytes.NewBuffer(in)

	// pkg: | 1B:0x38 | 4B:len | payload | 1B: 0x49 |
	var (
		pkgStx int8
		pkgLen int32
		pkgEtx int8
	)
	// stx
	if err := binary.Read(b, binary.BigEndian, &pkgStx); err != nil {
		return nil, 0, err
	}
	if pkgStx != 0x38 {
		return nil, 0, errors.CodecReadInvalid
	}
	// len
	if err := binary.Read(b, binary.BigEndian, &pkgLen); err != nil {
		return nil, 0, err
	}
	if pkgLen > maxDefaultPkgSize {
		return nil, 0, errors.CodecRead2BigErr
	}

	totalLen := int(1 + 4 + pkgLen + 1)
	if len(in) < totalLen {
		return nil, 0, errors.CodecReadIncompleteErr
	}

	// payload
	payload := make([]byte, pkgLen)
	if err := binary.Read(b, binary.BigEndian, payload); err != nil {
		return nil, 0, errors.CodecReadErr
	}
	// etx
	if err := binary.Read(b, binary.BigEndian, &pkgEtx); err != nil {
		return nil, 0, err
	}
	if pkgEtx != 0x49 {
		return nil, 0, errors.CodecReadInvalid
	}
	response := &Response{}
	if err := proto.Unmarshal(payload, response); err != nil {
		return nil, 0, err
	}

	return response, totalLen, nil
}
