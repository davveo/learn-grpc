package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/davveo/learn-grpc/proto/echo"

	"google.golang.org/grpc"
)

const addr = "localhost:50051"

type ecServer struct {
	echo.UnimplementedEchoServer
	addr string
}

func (s *ecServer) Simple(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	return &echo.Response{Message: fmt.Sprintf("%s (from %s)", req.Message, s.addr)}, nil
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &ecServer{addr: addr})
	log.Printf("serving on %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
