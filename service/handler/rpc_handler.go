package handler

// TODO: Flatten methods and functions.
// Either make function fulfill method interface or implement rpcMethod logic for
// functions.
// Can handler handle multiple functions? probably no
// TODO: Implement Call stack service.Call calls handler.Call calls method.Call
// TODO: Create Request interface

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/blushft/meld/service/handler/function"
	"github.com/blushft/meld/service/handler/method"
	"github.com/blushft/meld/utility"
)

type rpcHandler struct {
	v         interface{}
	methods   []method.Method
	functions []function.Function
	opts      *HandlerOptions
}

func NewHandler(v interface{}, opts ...HandlerOption) Handler {
	return newRPCHandler(v)
}

func newRPCHandler(v interface{}, opts ...HandlerOption) Handler {

	options := &HandlerOptions{
		Meta: make(map[string]map[string]string),
		Type: "rpc",
	}
	for _, o := range opts {
		o(options)
	}

	attr := utility.Attr(v)
	spew.Dump(attr)
	hType := reflect.TypeOf(v)
	// handler := reflect.ValueOf(v)

	methods := make([]method.Method, 0)
	functions := make([]function.Function, 0)
	var err error

	switch attr["kind"] {
	case "struct":
		methods, err = extractMethods(hType)
		if err != nil {
			log.Println(err)
		}
	case "func":
		fn, err := extractFunc(v)
		if err != nil {
			log.Println(err)
		}
		functions = append(functions, fn)
	}

	if len(options.Name) == 0 {
		options.Name = attr["name"]
	}
	h := &rpcHandler{
		functions: functions,
		methods:   methods,
		opts:      options,
		v:         v,
	}

	return h
}

func extractMethods(typ reflect.Type) ([]method.Method, error) {
	methods := make([]method.Method, 0)
	errs := make([]string, 0)
	var retErr error
	for m := 0; m < typ.NumMethod(); m++ {
		e, err := method.NewMethod(typ.Method(m))
		if err != nil {
			errs = append(errs, err.Error())
		}
		if e != nil {
			methods = append(methods, e)
		}

	}
	if len(errs) > 0 {
		retErr = errors.New(strings.Join(errs, "; "))
	}
	return methods, retErr
}

func extractFunc(v interface{}) (function.Function, error) {
	fn, err := function.NewFunction(v)
	if err != nil {
		return nil, err
	}
	return fn, nil
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
	cerr := m.Call(in)
	err, ok := cerr[0].Interface().(error)
	if !ok || err != nil {
		return err
	}

	return nil
}

func (r *rpcHandler) Meta() map[string]map[string]string {
	return r.opts.Meta
}
