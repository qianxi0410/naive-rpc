package naiverpc

type options struct {
	conf string
}

type Option func(*options)

func WithConf(fpath string) Option {
	return func(o *options) {
		o.conf = fpath
	}
}
