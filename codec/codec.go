package codec

import "sync"

// encode * decode
// client and server side should impl the codec
type Codec interface {
	// codec name
	Name() string

	// Encode encode 'pkg' into []byte, here 'pkg' must be []byte,
	// which is marshaled and compressed data.
	Encode(pkg interface{}) (data []byte, err error)

	// Decode decode []byte, return decoded body data and length,
	// which is marshaled and compressed data.
	Decode(data []byte) (req interface{}, n int, err error)
}

var (
	codecMutex = sync.RWMutex{}
	codecs     = map[string]codec{}
)

type codec struct {
	// name
	name string
	// client side's codec
	clientCodec Codec
	// server side's codec
	serverCodec Codec
}

// register the client and server codec to framework
func RegisterCodec(proto string, server, client Codec) {
	codecMutex.Lock()
	defer codecMutex.Unlock()

	codecs[proto] = codec{
		name:        proto,
		serverCodec: server,
		clientCodec: client,
	}
}

// retrun server side's codec
func ServerCodec(proto string) Codec {
	codecMutex.RLock()
	defer codecMutex.RUnlock()

	return codecs[proto].serverCodec
}

// return client side's codec
func ClientCodec(proto string) Codec {
	codecMutex.RLock()
	defer codecMutex.RUnlock()

	return codecs[proto].clientCodec
}
