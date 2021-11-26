package server

import "github.com/qianxi0410/naive-rpc/router"

type options struct {
	Router *router.Router
}

type Option func(*options)

func WithRouter(r *router.Router) Option {
	return func(o *options) {
		o.Router = r
	}
}
