package server

import (
	"context"

	"github.com/google/uuid"

	"github.com/blushft/meld/context/metadata"

	"github.com/blushft/meld/provider"
	"github.com/blushft/meld/server/encoding"
	"github.com/blushft/meld/server/transport"
)

type Options struct {
	Name string
	ID   string

	Address string
	Host    string
	Port    string

	Context context.Context

	Providers map[string]provider.Provider
	Encoders  map[string]encoding.Register
	Transport transport.Transport
	Meta      metadata.Metadata
}

type Option func(*Options)

func NewOptions(opt ...Option) Options {
	opts := Options{
		ID:        uuid.New().String(),
		Port:      "0",
		Context:   context.Background(),
		Providers: make(map[string]provider.Provider),
		Encoders:  make(map[string]encoding.Register),
		Meta:      metadata.NewMetadata(),
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

func Port(p string) Option {
	return func(o *Options) {
		o.Port = p
	}
}
