package service

import (
	"encoding/json"
	"reflect"
	"sync"

	"github.com/blushft/meld/service/method"
)

type Service interface {
	Name() string
	Options() *Options
	Register(...Option)
	NewMethod(method.Method, ...method.MethodOption) error
	Methods() map[string]method.Method
	Usage() string
}

type service struct {
	name string
	opts *Options

	sync.RWMutex
	methods map[string]method.Method
}

func NewService(h interface{}, opts ...Option) Service {
	return newService(h, opts...)
}

func newService(h interface{}, opts ...Option) Service {

	hType := reflect.TypeOf(h)
	handler := reflect.ValueOf(h)
	n := reflect.Indirect(handler).Type().Name()

	opts = append(opts, Name(n))
	options := newOptions(opts...)

	methods := make(map[string]method.Method)

	for m := 0; m < hType.NumMethod(); m++ {
		if e, _ := method.NewMethod(hType.Method(m)); e != nil {
			mName := n + "." + e.Name()

			for k, v := range options.Meta[mName] {
				e.Options().Meta[k] = v
			}

			methods[mName] = e
		}
	}

	return &service{
		name:    n,
		opts:    options,
		methods: methods,
	}
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

func (s *service) Register(opts ...Option) {
	s.Lock()
	defer s.Unlock()
	for _, o := range opts {
		o(s.opts)
	}
}

func (s *service) NewMethod(m method.Method, opts ...method.MethodOption) error {
	return nil
}

func (s *service) Methods() map[string]method.Method {
	return s.methods
}

func (s *service) Usage() string {
	methods := make(map[string]interface{})
	for _, m := range s.methods {
		methods[m.Name()] = map[string]interface{}{
			"request":  m.Request(),
			"response": m.Response(),
		}
	}

	usage := map[string]interface{}{
		"name":    s.name,
		"version": s.opts.Version.String(),
		"tags":    s.opts.Tags,
		"labels":  s.opts.Labels,
		"methods": methods,
	}

	outb, _ := json.MarshalIndent(usage, "", "  ")
	return string(outb)
}
