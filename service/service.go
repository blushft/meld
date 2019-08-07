package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

type Service interface {
	Name() string
	Options() *Options
	Configure(...Option)
	Handle(Handler, ...HandlerOption) error
	Handler(string) Handler
	Handlers() []string
	Usage() string
	Call(req Request, resp interface{}) error
}

func NewService(h interface{}, opts ...Option) Service {
	return newService(h, opts...)
}

type service struct {
	opts *Options

	sync.RWMutex
	handlers map[string]Handler
}

func newService(v interface{}, opts ...Option) Service {

	options := newOptions(opts...)

	if options.Name == "" {
		h := reflect.ValueOf(v)
		n := reflect.Indirect(h).Type().Name()
		options.Name = n
	}

	handlers := make(map[string]Handler)

	s := &service{
		handlers: handlers,
		opts:     options,
	}

	s.addHandler(v)

	return s

}

func (s *service) Name() string {
	s.RLock()
	name := s.opts.Name
	s.RUnlock()
	return name
}

func (s *service) Options() *Options {
	s.RLock()
	options := s.opts
	s.RUnlock()
	return options
}

func (s *service) Configure(opts ...Option) {
	s.Lock()
	defer s.Unlock()
	for _, o := range opts {
		o(s.opts)
	}
}

func (s *service) Handle(h Handler, opts ...HandlerOption) error {
	s.handlers[h.Name()] = h
	return nil
}

func (s *service) addHandler(v interface{}, opts ...HandlerOption) error {
	h := reflect.ValueOf(v)
	n := reflect.Indirect(h).Type().Name()
	th := NewHandler(v)

	s.handlers[n] = th
	return nil
}

func (s *service) Handler(n string) Handler {
	if h, ok := s.handlers[n]; ok {
		return h
	}

	return nil
}

func (s *service) Handlers() []string {
	ret := make([]string, 0)
	for k := range s.handlers {
		ret = append(ret, k)
	}
	return ret
}

func (s *service) Usage() string {
	handlers := make(map[string]map[string]interface{})
	for _, h := range s.handlers {
		methods := make(map[string]interface{})
		for _, m := range h.Methods() {
			methods[m.Name()] = map[string]interface{}{
				"request":  m.Request(),
				"response": m.Response(),
			}
		}
		handlers[h.Name()] = map[string]interface{}{"methods": methods}
	}

	usage := map[string]interface{}{
		"name":     s.Name(),
		"version":  s.opts.Version.String(),
		"tags":     s.opts.Tags,
		"labels":   s.opts.Labels,
		"handlers": handlers,
	}

	outb, _ := json.MarshalIndent(usage, "", "  ")
	return string(outb)
}

func (s *service) Call(req Request, resp interface{}) error {

	h, ok := s.handlers[req.Handler()]
	if !ok {
		return fmt.Errorf("no handler registered for %s", req.Handler())
	}

	if err := h.Call(req, resp); err != nil {
		return err
	}

	return nil
}
