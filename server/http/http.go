package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"

	"github.com/blushft/meld/server"
	"github.com/blushft/meld/service"
)

type httpServer struct {
	logger  server.Logger
	options server.Options

	svcs    map[string]service.Service
	router  server.Router
	routers map[string]server.Router

	mux      *http.ServeMux
	Server   *http.Server
	Listener net.Listener
}

func New(opts ...server.Option) (server.Server, error) {
	return new(opts...)
}

func new(opts ...server.Option) (server.Server, error) {
	s := &httpServer{
		logger:  server.NewLogger(),
		options: server.NewOptions(),
		svcs:    make(map[string]service.Service),
		routers: make(map[string]server.Router),
		mux:     http.NewServeMux(),
	}

	for _, o := range opts {
		o(&s.options)
	}

	return s, nil
}

func (s *httpServer) Name() string {
	return "meld_server"
}

func (s *httpServer) Options() server.Options {
	return s.options
}

func (s *httpServer) Register(svcs ...service.Service) error {
	for _, svc := range svcs {
		s.svcs[s.Name()] = svc
		for _, hdl := range svc.Handlers() {
			s.makeHandlers(svc, svc.Handler(hdl))
		}
	}

	return nil
}

func (s *httpServer) Configure(opts ...server.Option) {
	for _, o := range opts {
		o(&s.options)
	}
}

func (s *httpServer) Routes() map[string]server.Route {
	return nil
}

func (s *httpServer) Start() error {
	addr := fmt.Sprintf("%s:%s", s.options.Host, s.options.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.logger.Logf("starting meld server on %s", l.Addr().String())

	return http.Serve(l, s.mux)
}

func (s *httpServer) Stop() error {
	return nil
}

func (s *httpServer) makeHandlers(svc service.Service, h service.Handler) {
	path := fmt.Sprintf("/%s/%s/", svc.Name(), h.Name())
	for _, m := range h.Methods() {
		mpath := path + m.Name()

		s.mux.HandleFunc(mpath, func(w http.ResponseWriter, r *http.Request) {
			var reqType, respType reflect.Type
			var req, resp reflect.Value
			var isVal, hasReq, hasResp bool

			if m.Request() != nil {
				hasReq = true
				reqType = m.Request().TypeOf()
				if reqType.Kind() == reflect.Ptr {
					req = reflect.New(reqType.Elem())
				} else {
					req = reflect.New(reqType)
					isVal = true
				}
			}

			if m.Response() != nil {
				hasResp = true
				respType = m.Response().TypeOf()
			}

			s.logger.Logf("hit on %s\n", mpath)
			reqOpts := []service.RequestOption{
				service.WithContext(r.Context()),
			}

			if hasReq {
				reqBody, _ := ioutil.ReadAll(r.Body)
				if err := json.Unmarshal(reqBody, req.Interface()); err != nil {
					s.logger.Log(err)
					http.Error(w, "invalid request", http.StatusBadRequest)
					return
				}
				if isVal {
					req = req.Elem()
				}
				reqOpts = append(reqOpts, service.WithBody("interface", req.Interface()))
			} else {
				req = reflect.ValueOf(nil)
			}

			hreq := h.NewRequest(
				m.Name(),
				reqOpts...,
			)

			if hasResp {
				if respType.Kind() == reflect.Ptr {
					resp = reflect.New(respType.Elem())
				} else {
					resp = reflect.New(respType)
				}
			} else {
				resp = reflect.ValueOf(nil)
			}

			if err := svc.Call(hreq, resp.Interface()); err != nil {
				s.logger.Log(err)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			if err := json.NewEncoder(w).Encode(resp.Interface()); err != nil {
				s.logger.Log(err)
				http.Error(w, "error writing response", http.StatusInternalServerError)
			}
		})
	}
}

func (s *httpServer) newMethodHandler(m service.Method) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("you reached %s", m.Name())))
	}
}
