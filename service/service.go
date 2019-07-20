package service

import (
	"context"

	"github.com/blushft/meld/service/handler"
)

type Service interface {
	Name() string
	Options() *Options
	Configure(...Option)
	Handle(handler.Handler, ...handler.HandlerOption) error
	Handler(string) handler.Handler
	Handlers() []string
	Usage() string
	Call(ctx context.Context, handler, method string, req, resp interface{}, opts ...CallOptions) error
}

func NewService(h interface{}, opts ...Option) Service {
	return newRPCService(h, opts...)
}
