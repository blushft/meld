package handler

import (
	"reflect"
)

type Handler interface {
	Name() string
	Request() *HandlerDef
	Response() *HandlerDef
	Meta() map[string]string
}

type HandlerDef struct {
	Name   string        `json:"name"`
	Type   string        `json:"type"`
	Values []*HandlerDef `json:"arguments"`
}

type HandlerInvocation func(name string, reqType, respType reflect.Type) (Handler, error)
