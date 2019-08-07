package service

import (
	"context"
	"reflect"
	"strings"

	"github.com/blang/semver"
)

type Method interface {
	Name() string
	Request() *MethodDef
	Response() *MethodDef
	Options() *MethodOptions
}

type MethodDef struct {
	Name   string       `json:"name"`
	Type   string       `json:"type"`
	Values []*MethodDef `json:"arguments"`

	t reflect.Type
}

func (m *MethodDef) TypeOf() reflect.Type {
	return m.t
}

type MethodOptions struct {
	Name    string
	Version semver.Version

	Labels map[string]string
	Tags   []string
}

type MethodOption func(*MethodOptions)

type MethodInvocation func(rm reflect.Method, opts ...MethodOption) (Method, error)

type method struct {
	name     string
	request  *MethodDef
	response *MethodDef
	opts     *MethodOptions
}

func NewMethod(rm reflect.Method, opts ...MethodOption) (Method, error) {
	return newMethod(rm, opts...)
}

func newMethod(rm reflect.Method, opts ...MethodOption) (Method, error) {
	options := &MethodOptions{
		Labels: make(map[string]string),
		Tags:   make([]string, 0),
	}

	for _, o := range opts {
		o(options)
	}

	m := &method{
		opts: options,
	}

	if err := m.extractMethod(rm); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *method) Name() string {
	return m.opts.Name
}

func (m *method) Request() *MethodDef {
	return m.request
}

func (m *method) Response() *MethodDef {
	return m.response
}

func (m *method) Options() *MethodOptions {
	return m.opts
}

func (m *method) Call(ctx context.Context, req interface{}, resp interface{}, opts ...MethodOption) error {
	return nil
}

func (m *method) extractMethod(rm reflect.Method) error {
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

	m.request = extractMethod(reqType, 0)
	m.response = extractMethod(respType, 0)
	if m.opts.Name == "" {
		m.opts.Name = rm.Name
	}
	return nil
}

func extractMethod(v reflect.Type, d int) *MethodDef {
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
		t:    v,
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			fval := extractMethod(f.Type, d+1)
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
		fval := extractMethod(v.Elem(), d+1)
		if fval != nil {
			val.Values = append(val.Values, fval)
		}
	}

	return val
}
