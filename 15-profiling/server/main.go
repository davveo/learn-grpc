package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/davveo/learn-grpc/proto/echo"

	"google.golang.org/grpc"
	profsvc "google.golang.org/grpc/profiling/service"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	echo.UnimplementedEchoServer
}

func (s *server) Simple(ctx context.Context, in *echo.Request) (*echo.Response, error) {
	fmt.Printf("Simple called with message %q\n", in.GetMessage())
	return &echo.Response{Message: in.Message}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &server{})

	// Register your grpc.Server with profiling.
	pc := &profsvc.ProfilingConfig{
		Server:          s,
		Enabled:         true,
		StreamStatsSize: 1024,
	}
	if err = profsvc.Init(pc); err != nil {
		fmt.Printf("error calling profsvc.Init: %v\n", err)
		return
	}

	s.Serve(lis)
}
