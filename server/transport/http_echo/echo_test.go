package echo

import (
	"reflect"
	"testing"

	greeter "github.com/blushft/meld/examples/greeter/service"
	"github.com/blushft/meld/service"

	"github.com/blushft/meld/server"
)

var (
	greeterSvc = service.NewService(&greeter.Greeter{}, []service.Option{
		service.Name("testsvc"),
		service.Namespace("meld.test"),
		service.WithTag("test", "testing", "greet"),
		service.WithLabel("release", "latest"),
		service.Version("0.1.0"),
	}...)
)

func init() {
	es := NewEchoServer(
		server.Address("5499"),
	)

	es.Start()
}

func TestEchoServer(t *testing.T) {
	type args struct {
		opts []server.Option
	}
	tests := []struct {
		name string
		args args
		want server.Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEchoServer(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEchoServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
