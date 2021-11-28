package pool

import "time"

type Options struct {
	MinIdle int
	// -1 represent no limit
	MaxIdle int
	// -1 represent no limit
	MaxActive int
	// whether wait if reach max conn
	Wait bool
	// idle conn timeout time
	IdleTimeout     time.Duration
	MaxConnLifetime time.Duration
	DialTimeout     time.Duration
}

type Option func(*Options)

// some helper funcs
func WithMinIdle(n int) Option {
	return func(o *Options) {
		o.MinIdle = n
	}
}

func WithMaxIdle(m int) Option {
	return func(o *Options) {
		o.MaxIdle = m
	}
}

func WithMaxActive(s int) Option {
	return func(o *Options) {
		o.MaxActive = s
	}
}

func WithWait(w bool) Option {
	return func(o *Options) {
		o.Wait = w
	}
}

func WithIdleTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = t
	}
}

func WithMaxConnLifetime(t time.Duration) Option {
	return func(o *Options) {
		o.MaxConnLifetime = t
	}
}

func WithDialTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = t
	}
}
