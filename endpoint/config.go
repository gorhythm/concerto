// Copyright 2024 The Concerto Contributors.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package endpoint

import (
	"slices"

	"github.com/go-kit/kit/endpoint"
)

// Config is a group of options for a endpoint Set.
type Config struct {
	middlewares []endpoint.Middleware
}

// ApplyMiddlewares applies the configured middlewares to the provided
// endpoint.
func (c Config) ApplyMiddlewares(e endpoint.Endpoint) endpoint.Endpoint {
	middlewares := slices.Clone(c.middlewares)
	slices.Reverse(c.middlewares)

	for i := len(middlewares) - 1; i >= 0; i-- {
		e = middlewares[i](e)
	}

	return e
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

// WithMiddlewares returns an [Option] that sets middlewares to the configured
// endpoint.
func WithMiddlewares(
	middlewares ...endpoint.Middleware,
) Option {
	return optionFunc(func(c Config) Config {
		c.middlewares = middlewares
		return c
	})
}
