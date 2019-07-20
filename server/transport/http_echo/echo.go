package echo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/blushft/meld/lifecycle"
	"github.com/blushft/meld/server"
	"github.com/blushft/meld/service"
	"github.com/labstack/echo"
)

func NewEchoServer(opts ...server.Option) server.Server {
	return newEchoServer(opts...)
}

type echoServer struct {
	e    *echo.Echo
	opts server.Options
	svcs []service.Service

	exit chan chan error
	sync.RWMutex

	wg sync.WaitGroup
	lifecycle.Lifecycle
}

func newEchoServer(opts ...server.Option) server.Server {
	options := server.NewOptions(opts...)
	options.Name = "Meld_HTTP_Echo_Server"
	return &echoServer{
		opts: options,
		e:    echo.New(),
		svcs: make([]service.Service, 0),
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

func (s *echoServer) Register(svcs ...service.Service) error {
	s.Lock()
	for _, svc := range svcs {
		s.svcs = append(s.svcs, svc)
	}
	return nil
}

func (s *echoServer) Configure(opts ...server.Option) {
	s.Lock()
	for _, opt := range opts {
		opt(&s.opts)
	}
	s.Unlock()
}

func (s *echoServer) Services() []service.Service {
	s.RLock()
	svcs := []service.Service{}
	for _, svc := range s.svcs {
		svcs = append(svcs, svc)
	}
	s.RUnlock()
	if len(svcs) > 0 {
		return svcs
	}
	return nil
}

func (s *echoServer) Start() error {
	smap := s.serviceMap()
	s.e.GET("/debug/servicemap", func(c echo.Context) error {
		return c.JSON(http.StatusOK, smap)
	})
	s.e.GET("/debug/usage", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, s.getUsage(), "  ")
	})

	s.e.HideBanner = true
	s.e.HidePort = true
	addr := fmt.Sprintf("%s:%s", s.opts.Address, s.opts.Port)
	return s.e.Start(addr)
}

func (s *echoServer) Stop() error {
	return nil
}

func (s *echoServer) serviceMap() map[string]map[string][]string {
	smap := make(map[string]map[string][]string)
	for _, svc := range s.svcs {
		hp := make(map[string][]string)
		for _, hn := range svc.Handlers() {
			mp := []string{}
			h := svc.Handler(hn)
			for _, m := range h.Methods() {
				mp = append(mp, m.Name())
				s.setupHandler(svc, hn, m.Name(), m.Request().TypeOf(), m.Response().TypeOf())
				hp[hn] = mp
			}
		}
		smap[svc.Name()] = hp
		s.e.GET(fmt.Sprintf("/%s/usage", svc.Name()), func(c echo.Context) error {
			return c.JSONBlob(http.StatusOK, []byte(svc.Usage()))
		})
	}

	return smap
}

func (s *echoServer) setupHandler(svc service.Service, hndl string, meth string, reqType, respType reflect.Type) {

	path := fmt.Sprintf("/%s/%s/%s", svc.Name(), hndl, meth)
	s.e.Add(
		"POST",
		path,
		func(c echo.Context) error {
			var req, resp reflect.Value
			var isVal bool
			if reqType.Kind() == reflect.Ptr {
				req = reflect.New(reqType.Elem())
			} else {
				req = reflect.New(reqType)
				isVal = true
			}

			reqBody, _ := ioutil.ReadAll(c.Request().Body)
			if err := json.Unmarshal(reqBody, req.Interface()); err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}

			if isVal {
				req = req.Elem()
			}

			isVal = false
			if respType.Kind() == reflect.Ptr {
				resp = reflect.New(respType.Elem())
			} else {
				resp = reflect.New(respType)
			}

			if err := svc.Call(c.Request().Context(), hndl, meth, req.Interface(), resp.Interface()); err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

			return c.JSON(http.StatusOK, resp.Interface())
		})
}

func (s *echoServer) getUsage() []map[string]interface{} {
	usage := make([]map[string]interface{}, 0)
	for _, svc := range s.svcs {
		var u map[string]interface{}
		if err := json.Unmarshal([]byte(svc.Usage()), &u); err != nil {
			log.Println(err)
		}
		usage = append(usage, u)
	}
	return usage
}
