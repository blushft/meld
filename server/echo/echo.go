package echo

import (
	"sync"

	"github.com/blushft/meld/service"

	"github.com/blushft/meld/server"
	"github.com/labstack/echo"
)

func NewEchoServer(opts ...server.Option) server.Server {
	return newEchoServer(opts...)
}

type echoServer struct {
	e    *echo.Echo
	opts server.Options

	exit chan chan error
	sync.RWMutex

	wg sync.WaitGroup
}

func newEchoServer(opts ...server.Option) server.Server {
	options := server.NewOptions(opts...)
	return &echoServer{
		opts: options,
		e:    echo.New(),
		exit: make(chan chan error),
	}
}

func (s *echoServer) Name() string {
	s.RLock()
	name := s.opts.Name
	s.RUnlock()
	return name
}

func (s *echoServer) Options() server.Options {
	s.RLock()
	opts := s.opts
	s.RUnlock()
	return opts
}

func (s *echoServer) Register(opts ...server.Option) error {
	s.Lock()
	for _, opt := range opts {
		opt(&s.opts)
	}
	s.Unlock()
	return nil
}

func (s *echoServer) Services() []service.Service {
	s.RLock()
	svcs := []service.Service{}
	for _, svc := range s.opts.Services {
		svcs = append(svcs, svc)
	}
	s.RUnlock()
	if len(svcs) > 0 {
		return svcs
	}
	return nil
}

func (s *echoServer) Start() error {
	return nil
}

func (s *echoServer) Stop() error {
	return nil
}
