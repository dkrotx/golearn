package main

import (
	"errors"
	"time"
)

var KeyNotFoundErr = errors.New("key not found")
var KeyExists = errors.New("key already exists")

type SetOptions struct {
	TTL time.Time
	OnlyIfNotExists bool
}

// use functional-params idiom

type Option func(opts *SetOptions)

func TTL(ttl time.Duration) Option {
	return func(opts *SetOptions) {
		opts.TTL = time.Now().Add(ttl)
	}
}

func IfNotExists() Option {
	return func(opts *SetOptions) {
		opts.OnlyIfNotExists = true
	}
}

func ApplyOptions(origin *SetOptions, opts...Option) {
	for _, opt := range opts {
		opt(origin)
	}
}

type KeyValue interface {
	Get(key string) (value string, err error)
	Set(key, value string, opts...Option) error
	Delete(key string) error
}
