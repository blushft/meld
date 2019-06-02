package server

import (
	"github.com/blushft/meld/service"
)

type Server interface {
	Name() string
	Options() Options
	Register(...Option) error
	Services() []service.Service
	Start() error
	Stop() error
}

type Handler interface {
	Name() string
	Handler() interface{}
	// Endpoints() []router.Endpoint
	Options() HandlerOptions
}
