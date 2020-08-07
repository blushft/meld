package provider

import (
	"github.com/blushft/meld/service"
)

type Provider interface {
	Name() string
	WrapHandler(h service.Handler) service.Handler
	WrapMethod(m service.Method) service.Method
	WrapService(s service.Service) service.Service
}
