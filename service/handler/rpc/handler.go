package rpc

import (
	"reflect"
	"strings"

	"github.com/blushft/meld/service/handler"
)

type rpcHandler struct {
	name     string
	request  *handler.HandlerDef
	response *handler.HandlerDef
	meta     map[string]string
}

func NewHandler(name string, reqType, respType reflect.Type) (handler.Handler, error) {
	return newHandler(name, reqType, respType)
}

func newHandler(name string, reqType, respType reflect.Type) (handler.Handler, error) {
	req := extractSig(reqType, 0)
	resp := extractSig(respType, 0)

	h := &rpcHandler{
		name:     name,
		request:  req,
		response: resp,
		meta:     make(map[string]string),
	}

	return h, nil
}

func (r *rpcHandler) Name() string {
	return r.name
}

func (r *rpcHandler) Request() *handler.HandlerDef {
	return r.request
}

func (r *rpcHandler) Response() *handler.HandlerDef {
	return r.response
}

func (r *rpcHandler) Meta() map[string]string {
	return r.meta
}

func extractSig(v reflect.Type, d int) *handler.HandlerDef {
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

	val := &handler.HandlerDef{
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
