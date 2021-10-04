package greeter

import (
	"context"
	"fmt"

	pb "github.com/blushft/meld/examples/grpc/proto"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	msg := fmt.Sprintf("Hello, %s", req.Name)

	return &pb.HelloResponse{
		Message: msg,
	}, nil
}
