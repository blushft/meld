package provider

import (
	"github.com/blushft/meld/service"
	"github.com/blushft/meld/service/handler"
	"github.com/blushft/meld/service/handler/method"
)

type Provider interface {
	Name() string
	WrapHandler(h handler.Handler) handler.Handler
	WrapMethod(m method.Method) method.Method
	WrapService(s service.Service) service.Service
}
