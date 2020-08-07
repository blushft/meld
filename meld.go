package meld

import (
	"github.com/blushft/meld/service"
)

// NewService is a convenience method for service.NewService()
func NewService(h interface{}, opts ...service.Option) service.Service {
	return service.NewService(h, opts...)
}
