package service

import (
	"reflect"
	"strings"

	"github.com/blang/semver"
)

type Function interface {
	Name() string
	Request() *FuncDef
	Response() *FuncDef
	Options() *FuncOptions
}

type FuncDef struct {
	Name   string     `json:"name"`
	Type   string     `json:"type"`
	Values []*FuncDef `json:"values"`

	t reflect.Type
}

func (f *FuncDef) TypeOf() reflect.Type {
	return f.t
}

type FuncOptions struct {
	Name    string
	Version semver.Version
	Labels  map[string]string
	Tags    []string
}

type FuncOption func(*FuncOptions)

type FunctionInvocation func(rf reflect.Type, opts ...FuncOption) (Function, error)

type funcHndlr struct {
	args *FuncDef
	ret  *FuncDef
	opts *FuncOptions

	caller     reflect.Value
	isVariadic bool
}

func NewFunction(rf interface{}, opts ...FuncOption) (Function, error) {
	return newFunc(rf, opts...)
}

func newFunc(rf interface{}, opts ...FuncOption) (Function, error) {
	options := &FuncOptions{
		Labels: make(map[string]string),
		Tags:   make([]string, 0),
	}

	for _, o := range opts {
		o(options)
	}

	f := &funcHndlr{
		opts:   options,
		caller: reflect.ValueOf(rf),
	}

	if err := f.extractFunc(rf); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *funcHndlr) Name() string {
	return f.opts.Name
}

func (f *funcHndlr) Request() *FuncDef {
	return f.args
}

func (f *funcHndlr) Response() *FuncDef {
	return f.ret
}

func (f *funcHndlr) Options() *FuncOptions {
	return f.opts
}

func (f *funcHndlr) extractFunc(v interface{}) error {
	ft := reflect.TypeOf(v)
	if ft.Kind() != reflect.Func {
		return nil
	}

	f.isVariadic = ft.IsVariadic()

	args := &FuncDef{
		Name:   "args",
		Type:   "-",
		Values: make([]*FuncDef, 0),
	}

	for i := 0; i < ft.NumIn(); i++ {
		arg := extractFunction(ft.In(i), 0)
		if arg != nil {
			args.Values = append(args.Values, arg)
		}
	}

	f.args = args

	rets := &FuncDef{
		Name:   "returns",
		Type:   "-",
		Values: make([]*FuncDef, 0),
	}

	for i := 0; i < ft.NumOut(); i++ {
		ret := extractFunction(ft.Out(i), 0)
		if ret != nil {
			rets.Values = append(rets.Values, ret)
		}
	}

	f.ret = rets

	return nil
}

func extractFunction(v reflect.Type, d int) *FuncDef {
	if v == nil {
		return nil
	}

	vName := v.Name()

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		vName = "*" + v.Name()
	}

	if v.Kind() == reflect.Interface {
		vName = v.PkgPath() + "#" + v.Name()
	}

	val := &FuncDef{
		Name: vName,
		Type: v.Name(),
		t:    v,
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			fval := extractFunction(f.Type, d+1)
			if fval == nil {
				continue
			}

			if tags := f.Tag.Get("json"); len(tags) > 0 {
				tp := strings.Split(tags, ",")
				if tp[0] == "-" || tp[0] == "omitempty" {
					continue
				}
				fval.Name = tp[0]
			} else {
				fval.Name = ""
			}

			if len(fval.Name) == 0 {
				fval.Name = f.Name
			}

			if len(fval.Name) == 0 {
				continue
			}

			val.Values = append(val.Values, fval)
		}
	case reflect.Slice:
		p := v.Elem()
		if p.Kind() == reflect.Ptr {
			p = p.Elem()
		}
		val.Type = "[]" + p.Name()
		if len(val.Name) == 0 {
			val.Name = val.Type
		}
		fval := extractFunction(v.Elem(), d+1)
		if fval != nil {
			val.Values = append(val.Values, fval)
		}
	}

	return val
}
