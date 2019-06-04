package handler

import (
	"context"
	"reflect"

	"github.com/davecgh/go-spew/spew"

	"github.com/blushft/meld/service/method"
)

type rpcHandler struct {
	v       interface{}
	methods []method.Method
	opts    *HandlerOptions
}

func NewHandler(v interface{}, opts ...HandlerOption) (Handler, error) {
	return newRPCHandler(v)
}

func newRPCHandler(v interface{}, opts ...HandlerOption) (Handler, error) {

	hType := reflect.TypeOf(v)
	handler := reflect.ValueOf(v)
	n := reflect.Indirect(handler).Type().Name()

	options := &HandlerOptions{
		Name: n,
		Meta: make(map[string]map[string]string),
		Type: "rpc",
	}
	for _, o := range opts {
		o(options)
	}

	methods := make([]method.Method, 0)
	for m := 0; m < hType.NumMethod(); m++ {
		if e, _ := method.NewMethod(hType.Method(m)); e != nil {
			mName := n + "." + e.Name()

			for k, v := range options.Meta[mName] {
				e.Options().Meta[k] = v
			}

			methods = append(methods, e)
		}
	}

	h := &rpcHandler{
		methods: methods,
		opts:    options,
		v:       v,
	}

	return h, nil
}

func (r *rpcHandler) Name() string {
	return r.opts.Name
}

func (r *rpcHandler) Options() *HandlerOptions {
	return r.opts
}
func (r *rpcHandler) Methods() []method.Method {
	return r.methods
}

func (r *rpcHandler) Call(ctx context.Context, method string, req interface{}, resp interface{}, opts ...HandlerOption) error {
	m := reflect.ValueOf(r.v).MethodByName(method)
	in := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(req),
		reflect.ValueOf(resp),
	}
	spew.Dump(in)

	mresp := m.Call(in)
	spew.Dump(mresp[0].Interface())
	spew.Dump(resp)
	return nil

}

func (r *rpcHandler) Meta() map[string]map[string]string {
	return r.opts.Meta
}
