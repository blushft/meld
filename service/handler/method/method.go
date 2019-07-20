package method

import (
	"reflect"

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
