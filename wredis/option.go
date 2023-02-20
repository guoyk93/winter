package wredis

import "github.com/redis/go-redis/v9"

type KeyType string

const (
	Default KeyType = "default"
)

// Option function modifying options
type Option func(opts *options)

// WithKey change key for injection
func WithKey(k string) Option {
	return func(opts *options) {
		opts.key = KeyType(k)
	}
}

// WithURL set env key for redis options loading
func WithURL(k string) Option {
	return func(opts *options) {
		opts.url = k
	}
}

// WithOptions set [redis.Options] directly
func WithOptions(rOpts *redis.Options) Option {
	return func(opts *options) {
		opts.opts = rOpts
	}
}

type options struct {
	key  KeyType
	opts *redis.Options
	url  string
}

func createOptions(opts ...Option) *options {
	opt := &options{
		key: Default,
	}
	for _, item := range opts {
		item(opt)
	}
	return opt
}
