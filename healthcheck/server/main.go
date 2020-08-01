package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/davveo/learn-grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port  = flag.Int("port", 50051, "the port to serve on")
	sleep = flag.Duration("sleep", time.Second*5, "duration between changes in health")

	system = "" // empty string represents the health of the system
)

type echoServer struct {
	echo.UnimplementedEchoServer
}

func (e *echoServer) Simple(ctx context.Context, req *echo.Request) (*echo.Response, error) {
	return &echo.Response{
		Message: fmt.Sprintf("hello from localhost:%d", *port),
	}, nil
}

var _ echo.EchoServer = &echoServer{}

func main() {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(s, healthcheck)
	echo.RegisterEchoServer(s, &echoServer{})

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING

		for {
			healthcheck.SetServingStatus(system, next)

			if next == healthpb.HealthCheckResponse_SERVING {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}

			time.Sleep(*sleep)
		}
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
