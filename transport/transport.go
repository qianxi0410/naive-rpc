package transport

type Transport interface {
	ListenAndServe() error
	Closed() <-chan struct{}
	Network() string
	Address() string
	Codec() string
}
