package server

import (
	"github.com/google/uuid"

	"github.com/blushft/meld/context/metadata"
)

type Options struct {
	Name string
	ID   string
	Host string
	Port string

	Meta metadata.Metadata
}

type Option func(*Options)

func NewOptions(opt ...Option) Options {
	opts := Options{
		ID:   uuid.New().String(),
		Port: "0",
		Meta: metadata.NewMetadata(),
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

func Host(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func Port(p string) Option {
	return func(o *Options) {
		o.Port = p
	}
}
