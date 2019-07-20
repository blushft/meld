package function

import (
	"reflect"
	"strings"
)

type rpcFunc struct {
	args *FuncDef
	ret  *FuncDef
	opts *FuncOptions

	caller     reflect.Value
	isVariadic bool
}

func NewFunction(rf interface{}, opts ...FuncOption) (Function, error) {
	return newRPCFunc(rf, opts...)
}

func newRPCFunc(rf interface{}, opts ...FuncOption) (Function, error) {
	options := &FuncOptions{
		Labels: make(map[string]string),
		Tags:   make([]string, 0),
	}

	for _, o := range opts {
		o(options)
	}

	f := &rpcFunc{
		opts:   options,
		caller: reflect.ValueOf(rf),
	}

	if err := f.extractFunc(rf); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *rpcFunc) Name() string {
	return f.opts.Name
}

func (f *rpcFunc) Request() *FuncDef {
	return f.args
}

func (f *rpcFunc) Response() *FuncDef {
	return f.ret
}

func (f *rpcFunc) Options() *FuncOptions {
	return f.opts
}

func (f *rpcFunc) extractFunc(v interface{}) error {
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
		arg := extractSig(ft.In(i), 0)
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
		ret := extractSig(ft.Out(i), 0)
		if ret != nil {
			rets.Values = append(rets.Values, ret)
		}
	}

	f.ret = rets

	return nil
}

func extractSig(v reflect.Type, d int) *FuncDef {
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
			fval := extractSig(f.Type, d+1)
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
		fval := extractSig(v.Elem(), d+1)
		if fval != nil {
			val.Values = append(val.Values, fval)
		}
	}

	return val
}
