package server

import (
	"context"

	"github.com/blushft/meld/common/metadata"

	"github.com/blushft/meld/service"

	"github.com/blushft/meld/common/codec"
	"github.com/blushft/meld/common/transport"
	"github.com/blushft/meld/provider"
)

type Options struct {
	Name    string
	Address string
	ID      string
	Version string

	Context context.Context

	Providers map[string]provider.Provider
	Codecs    map[string]codec.Register
	Services  map[string]service.Service
	Transport transport.Transport
	Meta      metadata.Metadata
}

type Option func(*Options)

func NewOptions(opt ...Option) Options {
	opts := Options{
		Providers: make(map[string]provider.Provider),
		Codecs:    make(map[string]codec.Register),
		Services:  make(map[string]service.Service),
		Meta:      metadata.NewMetadata(),
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

func Service(s service.Service) Option {
	return func(o *Options) {
		o.Services[s.Name()] = s
	}
}
