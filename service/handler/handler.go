package handler

import (
	"context"

	"github.com/blushft/meld/service/handler/method"
)

type Handler interface {
	Name() string
	Options() *HandlerOptions
	Methods() []method.Method
	Call(ctx context.Context, method string, req interface{}, resp interface{}, opts ...HandlerOption) error
	Meta() map[string]map[string]string
}

type HandlerOptions struct {
	Name string
	Meta map[string]map[string]string
	Type string
}

type HandlerOption func(*HandlerOptions)

type HandlerInvocation func(v interface{}, opts ...HandlerOption) Handler

func Name(n string) HandlerOption {
	return func(o *HandlerOptions) {
		o.Name = n
	}
}
