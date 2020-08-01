package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	echo "github.com/davveo/learn-grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var port = flag.Int("port", 50052, "port number")

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

type server struct {
	echo.UnimplementedEchoServer
}

func (s *server) Simple(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	return &echo.Response{Message: req.Message}, nil
}

func main() {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)

	address := fmt.Sprintf(":%v", *port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	echo.RegisterEchoServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
