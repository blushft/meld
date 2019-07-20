package main

import (
	"context"
	"log"

	greeter "github.com/blushft/meld/examples/greeter/service"
	echo "github.com/blushft/meld/server/transport/http_echo"
	"github.com/blushft/meld/service"
	"github.com/blushft/meld/service/handler"

	"github.com/blushft/meld/server"
)

type StatusCheck struct{}

func (s *StatusCheck) Check(ctx context.Context, req struct{}, resp *string, opts ...handler.HandlerOption) error {
	*resp = "ok"
	return nil
}

var (
	es         server.Server
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

	greeterSvc.Handle(handler.NewHandler(&StatusCheck{}))
}

func main() {
	es = echo.NewEchoServer(
		server.Port("5499"),
	)

	es.Register(greeterSvc)

	if err := es.Start(); err != nil {
		log.Fatal(err)
	}
}
