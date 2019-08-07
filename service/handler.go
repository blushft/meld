package service

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/blushft/meld/utility"
)

type Handler interface {
	Name() string
	Options() *HandlerOptions
	Methods() []Method
	Call(req Request, resp interface{}, opts ...HandlerOption) error
	NewRequest(string, ...RequestOption) Request
	Meta() map[string]map[string]string
}

type HandlerOptions struct {
	Name string
	Meta map[string]map[string]string
	Type string
}

type HandlerOption func(*HandlerOptions)

type HandlerInvocation func(v interface{}, opts ...HandlerOption) Handler

func HandlerName(n string) HandlerOption {
	return func(o *HandlerOptions) {
		o.Name = n
	}
}

type handler struct {
	v         interface{}
	methods   []Method
	functions []Function
	opts      *HandlerOptions
}

func NewHandler(v interface{}, opts ...HandlerOption) Handler {
	return newHandler(v)
}

func newHandler(v interface{}, opts ...HandlerOption) Handler {

	options := &HandlerOptions{
		Meta: make(map[string]map[string]string),
		Type: "rpc",
	}
	for _, o := range opts {
		o(options)
	}

	attr := utility.Attr(v)
	//spew.Dump(attr)
	hType := reflect.TypeOf(v)
	// handler := reflect.ValueOf(v)

	methods := make([]Method, 0)
	functions := make([]Function, 0)
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
	h := &handler{
		functions: functions,
		methods:   methods,
		opts:      options,
		v:         v,
	}

	return h
}

func extractMethods(typ reflect.Type) ([]Method, error) {
	methods := make([]Method, 0)
	errs := make([]string, 0)
	var retErr error
	for m := 0; m < typ.NumMethod(); m++ {
		e, err := NewMethod(typ.Method(m))
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

func extractFunc(v interface{}) (Function, error) {
	fn, err := NewFunction(v)
	if err != nil {
		return nil, err
	}
	return fn, nil
}

func (r *handler) Name() string {
	return r.opts.Name
}

func (r *handler) Options() *HandlerOptions {
	return r.opts
}

func (r *handler) Methods() []Method {
	return r.methods
}

func (r *handler) NewRequest(method string, opts ...RequestOption) Request {
	fnd := false
	for _, m := range r.methods {
		if m.Name() == method {
			fnd = true
			continue
		}
	}

	if !fnd {
		return nil
	}

	reqOpts := RequestOptions{
		Handler: r.Name(),
		Method:  method,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&reqOpts)
	}

	return newOptRequest(reqOpts)
}

func (r *handler) Call(req Request, resp interface{}, opts ...HandlerOption) error {
	m := reflect.ValueOf(r.v).MethodByName(req.Method())
	in := []reflect.Value{
		reflect.ValueOf(req.Context()),
	}
	type noarg struct{}
	if req.Body() != nil {
		in = append(in, reflect.ValueOf(req.Body().Interface()))
	} else {
		in = append(in, reflect.ValueOf(noarg{}))
	}
	if resp != nil {
		in = append(in, reflect.ValueOf(resp))
	} else {
		in = append(in, reflect.ValueOf(noarg{}))
	}

	cerr := m.Call(in)
	err, ok := cerr[0].Interface().(error)
	if !ok || err != nil {
		return err
	}

	return nil
}

func (r *handler) Meta() map[string]map[string]string {
	return r.opts.Meta
}
