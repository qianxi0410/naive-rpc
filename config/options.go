package config

type Option func(*options)

type options struct {
	fp       string
	reload   bool
	decoder  Decoder
	provider Provider
}

func WithReload(v bool) Option {
	return func(o *options) {
		o.reload = v
	}
}

func WithDecoder(v Decoder) Option {
	return func(o *options) {
		o.decoder = v
	}
}

func WithProvider(v Provider) Option {
	return func(o *options) {
		o.provider = v
	}
}
