package main

import (
	"context"
	"flag"
	"fmt"
	echo "github.com/davveo/learn-grpc/pb"
	"github.com/davveo/learn-grpc/utils"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "the port to serve on")
	streamingCount = 10
)

type service struct {
	pb echo.UnimplementedEchoServer
}

func (s *service) ClientStream(stream echo.Echo_ClientStreamServer) error {
	panic("implement me")
}

func (s *service) ServerStream(request *echo.Request, stream echo.Echo_ServerStreamServer) error {
	panic("implement me")
}

func (s *service) DoubleStream(stream echo.Echo_DoubleStreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending echo\n", in)
		if err := stream.Send(&echo.Response{Message: in.Message}); err != nil {
			return err
		}
	}
}

func (s *service) Simple(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	return &echo.Response{Message: req.Message}, nil
}

func main()  {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)

	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &service{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	utils.CheckError(err)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
