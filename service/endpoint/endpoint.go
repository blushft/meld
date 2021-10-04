package endpoint

import (
	"errors"
	"reflect"
	"strings"

	"github.com/blushft/meld/pkg/reflectutil"
)

type Endpoint interface {
	Name() string
	Request() *Descriptor
	Response() *Descriptor
	Options() Options
	Call(Request, interface{}) error
}

type Descriptor struct {
	Name   string
	Type   string
	Values []*Descriptor

	t reflect.Type
}

type Options struct {
	Labels map[string]string
}

type Option interface {
	Apply(Endpoint)
}

type endpoint struct {
	name     string
	request  *Descriptor
	response *Descriptor
	opts     Options

	v interface{}
}

func (e *endpoint) Name() string {
	return e.name
}

func (e *endpoint) Request() *Descriptor {
	return e.request
}

func (e *endpoint) Response() *Descriptor {
	return e.response
}

func (e *endpoint) Options() Options {
	return e.opts
}

func (e *endpoint) Call(req Request, resp interface{}) error {
	m := reflect.ValueOf(e.v).MethodByName(e.name)

	in := []reflect.Value{
		reflect.ValueOf(req.Context()),
	}

	type noarg struct{}
	if req.Body() != nil {
		in = append(in, reflect.ValueOf(req.Body()))
	} else {
		in = append(in, reflect.ValueOf(noarg{}))
	}

	if resp != nil {
		in = append(in, reflect.ValueOf(resp))
	} else {
		in = append(in, reflect.ValueOf(noarg{}))
	}

	ret := m.Call(in)
	if len(ret) == 0 && resp != nil {
		return errors.New("call did not return expected result")
	}

	cerr, ok := ret[0].Interface().(error)
	if !ok || cerr != nil {
		return cerr
	}

	return nil
}

func Extract(v interface{}) ([]Endpoint, error) {
	var res []Endpoint
	var err error

	attr := reflectutil.Attr(v)

	switch attr.Kind {
	case "struct":
		res, err = extractEndpoints(v)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid endpoint extraction type")
	}

	return res, nil
}

func extractEndpoints(v interface{}) ([]Endpoint, error) {
	var eps []Endpoint

	typ := reflect.TypeOf(v)
	for m := 0; m < typ.NumMethod(); m++ {
		ep := &endpoint{
			name: typ.Method(m).Name,
			v:    v,
		}

		if err := ep.extractMethod(typ.Method(m)); err != nil {
			return nil, err
		}

		if ep != nil {
			eps = append(eps, ep)
		}
	}

	return eps, nil
}

func (ep *endpoint) extractMethod(m reflect.Method) error {
	if m.PkgPath != "" {
		return nil
	}

	var reqt, respt reflect.Type
	mt := m.Type

	switch mt.NumIn() {
	case 3:
		reqt = mt.In(1)
		respt = mt.In(2)
	case 4:
		reqt = mt.In(2)
		respt = mt.In(3)
	default:
		return nil
	}

	ep.request = extractMethod(reqt, 0)
	ep.response = extractMethod(respt, 0)

	return nil
}

func extractMethod(v reflect.Type, d int) *Descriptor {
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

	val := &Descriptor{
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
