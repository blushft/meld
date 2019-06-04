package method

import (
	"context"
	"reflect"
	"strings"
)

type rpcMethod struct {
	name     string
	request  *MethodDef
	response *MethodDef
	opts     *MethodOptions
}

func NewMethod(rm reflect.Method, opts ...MethodOption) (Method, error) {
	return newRPCMethod(rm, opts...)
}

func newRPCMethod(rm reflect.Method, opts ...MethodOption) (Method, error) {
	options := &MethodOptions{
		Meta: make(map[string]string),
	}

	for _, o := range opts {
		o(options)
	}

	m := &rpcMethod{}

	if err := m.extractMethod(rm); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *rpcMethod) Name() string {
	return m.name
}

func (m *rpcMethod) Request() *MethodDef {
	return m.request
}

func (m *rpcMethod) Response() *MethodDef {
	return m.response
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

	m.request = extractSig(reqType, 0)
	m.response = extractSig(respType, 0)

	m.name = rm.Name
	return nil
}

func extractSig(v reflect.Type, d int) *MethodDef {
	if d == 3 {
		return nil
	}
	if v == nil {
		return nil
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if len(v.Name()) == 0 {
		return nil
	}

	val := &MethodDef{
		Name: v.Name(),
		Type: v.Name(),
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
		fval := extractSig(v.Elem(), d+1)
		if fval != nil {
			val.Values = append(val.Values, fval)
		}
	}

	return val
}
