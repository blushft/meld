package method

import (
	"context"
	"reflect"

	"github.com/blushft/meld/service/handler"
)

type Method interface {
	Name() string
	Request() *handler.HandlerDef
	Response() *handler.HandlerDef
	Metadata() map[string]string
	// Handle(h Handler) error
	Options() *MethodOptions
	Call(ctx context.Context, req interface{}, resp interface{}, opts ...MethodOption) error
}

type MethodOptions struct {
	Meta map[string]string
}

type MethodOption func(*MethodOptions)

type MethodInvocation func(rm reflect.Method, opts ...MethodOption) (Method, error)
