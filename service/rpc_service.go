package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/blushft/meld/service/handler"
)

type rpcService struct {
	opts *Options

	sync.RWMutex
	handlers map[string]handler.Handler
}

func newRPCService(v interface{}, opts ...Option) Service {

	options := newOptions(opts...)

	if options.Name == "" {
		h := reflect.ValueOf(v)
		n := reflect.Indirect(h).Type().Name()
		options.Name = n
	}

	handlers := make(map[string]handler.Handler)

	s := &rpcService{
		handlers: handlers,
		opts:     options,
	}

	s.addHandler(v)

	return s

}

func (s *rpcService) Name() string {
	s.RLock()
	name := s.opts.Name
	s.RUnlock()
	return name
}

func (s *rpcService) Options() *Options {
	s.RLock()
	options := s.opts
	s.RUnlock()
	return options
}

func (s *rpcService) Configure(opts ...Option) {
	s.Lock()
	defer s.Unlock()
	for _, o := range opts {
		o(s.opts)
	}
}

func (s *rpcService) Handle(h handler.Handler, opts ...handler.HandlerOption) error {
	s.handlers[h.Name()] = h
	return nil
}

func (s *rpcService) addHandler(v interface{}, opts ...handler.HandlerOption) error {
	h := reflect.ValueOf(v)
	n := reflect.Indirect(h).Type().Name()
	th := handler.NewHandler(v)

	s.handlers[n] = th
	return nil
}

func (s *rpcService) Handler(n string) handler.Handler {
	if h, ok := s.handlers[n]; ok {
		return h
	}

	return nil
}

func (s *rpcService) Handlers() []string {
	ret := make([]string, 0)
	for k := range s.handlers {
		ret = append(ret, k)
	}
	return ret
}

func (s *rpcService) Usage() string {
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

func (s *rpcService) Call(ctx context.Context, handler, method string, req, resp interface{}, opts ...CallOptions) error {

	h, ok := s.handlers[handler]
	if !ok {
		return fmt.Errorf("no handler registered for %s", handler)
	}

	if err := h.Call(ctx, method, req, resp); err != nil {
		return err
	}

	return nil
}
