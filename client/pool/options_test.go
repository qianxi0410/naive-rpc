package pool_test

import (
	"testing"
	"time"

	"github.com/qianxi0410/naive-rpc/client/pool"
	"github.com/stretchr/testify/assert"
)

func TestWithOptions(t *testing.T) {
	opts := &pool.Options{}
	pool.WithMinIdle(1)(opts)
	pool.WithMaxIdle(2)(opts)
	pool.WithMaxActive(10)(opts)
	pool.WithIdleTimeout(time.Second)(opts)
	pool.WithDialTimeout(time.Second)(opts)
	pool.WithMaxConnLifetime(time.Second * 60)(opts)
	pool.WithWait(true)(opts)

	assert.Equal(t, opts.MinIdle, 1)
	assert.Equal(t, opts.MaxIdle, 2)
	assert.Equal(t, opts.MaxActive, 10)
	assert.Equal(t, opts.IdleTimeout, time.Second)
	assert.Equal(t, opts.DialTimeout, time.Second)
	assert.Equal(t, opts.MaxConnLifetime, 60*time.Second)
	assert.Equal(t, opts.Wait, true)
}
