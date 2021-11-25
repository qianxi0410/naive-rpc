package naiverpc

type option struct {
	conf string
}

type Option func(*option)

func WithConf(fpath string) Option {
	return func(o *option) {
		o.conf = fpath
	}
}
