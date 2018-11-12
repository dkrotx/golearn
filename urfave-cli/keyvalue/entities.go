package main

import (
	"errors"
	"time"
)

var KeyNotFoundErr = errors.New("key not found")
var KeyExists = errors.New("key already exists")

type SetOptions struct {
	TTL time.Duration
}

// use functional-params idiom

type Option func(opts *SetOptions)

func TTL(ttl time.Duration) Option {
	return func(opts *SetOptions) {
		opts.TTL = ttl
	}
}

func CombineOptions(opts...Option) (res SetOptions) {
	for _, opt := range opts {
		opt(&res)
	}
	return
}

type KeyValue interface {
	Get(key string) (value string, err error)
	Set(key, value string, opts...Option) error
	Delete(key string) error
}
