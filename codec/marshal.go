package codec

import (
	"sync"

	"google.golang.org/protobuf/proto"
)

type MarshalType int

const (
	MarshalTypePB = MarshalType(iota)
	MarshalTypeJSON
)

var (
	marshalersMutex sync.RWMutex
	marshalers      = map[MarshalType]Marshaler{}
)

// Marshaler defines how to serialize/deserialize
type Marshaler interface {
	// Marshal serializes 'pkg' to []byte, return the data,
	// return error if any error encountered.
	Marshal(pkg interface{}) (data []byte, err error)

	// Unmarhshal deserializes []byte to pkg, which must be a pointer,
	// return error if any error encountered.
	Unmarshal(data []byte, pkg interface{}) error
}

// RegisterMarshaler registers a new marshaler implemention of specific MarshalType,
// this is concurrent-safe.
func RegisterMarshaler(typ MarshalType, marshaler Marshaler) {
	marshalersMutex.Lock()
	defer marshalersMutex.Unlock()
	marshalers[typ] = marshaler

}

// DeregisterMarshaler deregisters the marshaler of specific MarshalType,
// this is concurrent-safe.
func DeregisterMarshaler(typ MarshalType) {
	marshalersMutex.Lock()
	defer marshalersMutex.Unlock()
	delete(marshalers, typ)
}

// PBMarshaler google protocolbuffers marshaling
type PBMarshaler struct{}

// Marshal marshal pkg in protocolbuffers manner
func (m *PBMarshaler) Marshal(pkg interface{}) (data []byte, err error) {
	msg, ok := pkg.(proto.Message)
	if !ok {
		return nil, MarshalInvalidPBErr
	}
	return proto.Marshal(msg)
}

// Unmarshal unmarshal data into pkg, which must be proto.Message
func (m *PBMarshaler) Unmarshal(data []byte, pkg interface{}) error {
	msg, ok := pkg.(proto.Message)
	if !ok {
		return MarshalInvalidPBErr
	}
	return proto.Unmarshal(data, msg)
}
