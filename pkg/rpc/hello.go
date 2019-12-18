package rpc

import (
	"context"

	pb "github.com/willzhao90/hellobackend/out"
)

type Server struct{}

func (s *Server) GetHello(ctx context.Context, in *pb.GetHelloRequest) (out *pb.GetHelloResponse, err error) {
	out = &pb.GetHelloResponse{
		Name: in.Name,
	}
	return out, nil
}

func NewServer() *Server {
	return &Server{}
}
