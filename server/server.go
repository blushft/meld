package server

import (
	"github.com/blushft/meld/service"
)

type Server interface {
	Name() string
	Options() Options
	Register(...service.Service) error
	Configure(...Option)
	Routes() map[string]Route
	Start() error
	Stop() error
}
