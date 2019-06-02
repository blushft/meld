package method

import (
	"context"
	"reflect"

	"github.com/blushft/meld/service/handler"
	"github.com/blushft/meld/service/handler/rpc"
)

type rpcMethod struct {
	name string

	opts *MethodOptions
	h    handler.Handler
}

func NewMethod(rm reflect.Method, opts ...MethodOption) (Method, error) {
	return newMethod(rm, opts...)
}

func newMethod(rm reflect.Method, opts ...MethodOption) (Method, error) {
	options := &MethodOptions{
		Meta: make(map[string]string),
	}

	for _, o := range opts {
		o(options)
	}

	m := &rpcMethod{
		opts: options,
	}

	if err := m.extractMethod(rm); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *rpcMethod) Name() string {
	return m.name
}

func (m *rpcMethod) Request() *handler.HandlerDef {
	return m.h.Request()
}

func (m *rpcMethod) Response() *handler.HandlerDef {
	return m.h.Response()
}

func (m *rpcMethod) Metadata() map[string]string {
	return m.opts.Meta
}

func (m *rpcMethod) Options() *MethodOptions {
	return m.opts
}

func (m *rpcMethod) Call(ctx context.Context, req interface{}, resp interface{}, opts ...MethodOption) error {
	return nil
}

func (m *rpcMethod) extractMethod(rm reflect.Method) error {
	if rm.PkgPath != "" {
		return nil
	}

	var reqType, respType reflect.Type
	mt := rm.Type

	switch mt.NumIn() {
	case 4:
		reqType = mt.In(1)
		respType = mt.In(2)
	case 5:
		reqType = mt.In(2)
		respType = mt.In(3)
	default:
		return nil
	}

	switch respType.Kind() {
	case reflect.Func:
		//do something
	}

	h, err := rpc.NewHandler(rm.Name, reqType, respType)
	if err != nil {
		return nil
	}

	m.name = rm.Name
	m.h = h

	return nil
}
