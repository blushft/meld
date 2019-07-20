package connector

import (
	"github.com/blushft/meld/options"
)

var (
	connectors = make(map[string]*Connector)
)

type Connector interface {
	Name() string
	Namespace() string
	Description() string

	Labels() map[string]string
	Tags() []string

	Options() options.Options
	Configure(...options.Setter)
	Reload(...options.Setter)

	Run() error
	IsRunning() bool
	Start() error
	Stop() error

	Event(interface{}) interface{}
}
