// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package endpoint

import (
	"github.com/go-kit/kit/endpoint"
)

// Config is a group of options for a endpoint Set.
type Config struct {
	middlewares []endpoint.Middleware
}

// Middleware returns a composed middlewares. Requests will traverse them in
// the order they're declared.
func (c Config) Middleware() endpoint.Middleware {
	if len(c.middlewares) == 0 {
		return nil
	}

	return endpoint.Chain(c.middlewares[0], c.middlewares[1:]...)
}

// NewConfig applies all the options to a returned [Config].
func NewConfig(opts ...Option) Config {
	cfg := Config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	return cfg
}

// A Option sets options to config.
type Option interface {
	apply(Config) Config
}

type optionFunc func(Config) Config

func (fn optionFunc) apply(cfg Config) Config {
	return fn(cfg)
}

// WithMiddlewares returns an [Option] that appends middlewares to the
// configured endpoint.
func WithMiddlewares(
	middlewares ...endpoint.Middleware,
) Option {
	return optionFunc(func(c Config) Config {
		c.middlewares = append(c.middlewares, middlewares...)
		return c
	})
}
