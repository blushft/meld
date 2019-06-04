package method

import (
	"reflect"
)

type Method interface {
	Name() string
	Request() *MethodDef
	Response() *MethodDef
	Metadata() map[string]string
	Options() *MethodOptions
}

type MethodDef struct {
	Name   string       `json:"name"`
	Type   string       `json:"type"`
	Values []*MethodDef `json:"arguments"`
}

type MethodOptions struct {
	Meta map[string]string
}

type MethodOption func(*MethodOptions)

type MethodInvocation func(rm reflect.Method, opts ...MethodOption) (Method, error)
