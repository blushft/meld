package server

import (
	"github.com/blushft/meld/service"
	"github.com/blushft/meld/service/handler"
)

type Server interface {
	Name() string
	Options() Options
	Register(...service.Service) error
	Configure(...Option)
	Endpoints() map[string]Endpoint
	Start() error
	Stop() error
}

type Listener interface {
	Name() string
	Handler() interface{}
	// Endpoints() []router.Endpoint
	Options() handler.HandlerOptions
}
