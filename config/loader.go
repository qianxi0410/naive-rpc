package config

import (
	"context"
)

//  load the config, it may internally uses Provider to read config
type Loader interface {
	Load(ctx context.Context, fp string, opts ...Option) (Config, error)
}

type loader struct {
	opts   options
	config config
}

func (r *loader) Load(ctx context.Context, fp string, opts ...Option) (Config, error) {
	o := options{
		fp:       fp,
		reload:   r.opts.reload,
		decoder:  r.opts.decoder,
		provider: r.opts.provider,
	}
	for _, opt := range opts {
		opt(&o)
	}

	if o.reload {
		r.reload(ctx, fp, o)
	}

	return r.load(ctx, fp, o)
}

func (r *loader) reload(ctx context.Context, fp string, opts options) error {
	ch, err := opts.provider.Watch(ctx, fp)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case v := <-ch:
				if v.typ != Update {
					continue
				}
				r.load(ctx, fp, opts)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *loader) decode(ctx context.Context, data []byte, opts options) (interface{}, error) {
	var cfg interface{}
	var err error

	switch v := opts.decoder.(type) {
	case *YAMLDecoder:
		cfg, err = v.Decode(data)
		// y := YAMLConfig{}
		// simpleyaml.NewYaml()
	default:
		panic("not support decode tyoe")
	}

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (r *loader) load(ctx context.Context, fp string, opts options) (Config, error) {
	data, err := opts.provider.Load(ctx, fp)
	if err != nil {
		return nil, err
	}

	cfg, err := r.decode(ctx, data, opts)
	if err != nil {
		return nil, err
	}

	r.config.value.Store(cfg)
	return &r.config, nil

}

func NewLoader(ctx context.Context, opts ...Option) (Loader, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	return &loader{opts: o}, nil
}
