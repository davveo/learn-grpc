package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/davveo/learn-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"net"
	"time"
)

var port = flag.Int("port", 50051, "the port to serve on")

const (
	timestampFormat = time.StampNano
	streamingCount  = 10
)

type server struct {
	echo.UnimplementedEchoServer
}

func (s *server) Simple(ctx context.Context, in *echo.Request) (*echo.Response, error) {
	fmt.Printf("--- Simple ---\n")
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		grpc.SetTrailer(ctx, trailer)
	}()

	// read from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnaryEcho: failed to get metadata")
	}
	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp from metadata:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}
	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	grpc.SendHeader(ctx, header)

	fmt.Printf("request received: %v, sending echo\n", in)

	return &echo.Response{Message: in.Message}, nil

}


func (s *server) ClientStream(echo.Echo_ClientStreamServer) error {
	return nil

}
func (s *server) ServerStream(*echo.Request, echo.Echo_ServerStreamServer) error {
	return nil
}
func (s *server) DoubleStream(echo.Echo_DoubleStreamServer) error {
	return nil
}





func main() {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)

	rand.Seed(time.Now().UnixNano())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	s := grpc.NewServer()
	echo.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}