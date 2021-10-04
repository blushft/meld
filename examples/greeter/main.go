package main

import (
	"context"
	"log"

	greeter "github.com/blushft/meld/examples/greeter/service"
	"github.com/blushft/meld/server/http"
	"github.com/blushft/meld/service"

	"github.com/blushft/meld/server"
)

type StatusCheck struct{}

func (s *StatusCheck) Check(ctx context.Context, req interface{}, resp *string, opts ...service.HandlerOption) error {
	*resp = "ok"
	return nil
}

var (
	greeterSvc service.Service
)

func init() {

	greeterSvc = service.NewService(&greeter.Greeter{}, []service.Option{
		service.Name("testsvc"),
		service.Namespace("meld.test"),
		service.WithTag("test", "testing", "greet"),
		service.WithLabel("release", "latest"),
		service.Version("0.1.0"),
	}...)

	greeterSvc.Handle(service.NewHandler(&StatusCheck{}))
}

func main() {
	s, _ := http.New(
		server.Port("5499"),
	)

	s.Register(greeterSvc)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
