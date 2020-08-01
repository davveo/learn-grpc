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
	profsvc "google.golang.org/grpc/profiling/service"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")
var profilingPort = flag.Int("profilingPort", 50052, "port to expose the profiling service on")

func setupClientProfiling() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *profilingPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return err
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()

	// Register this grpc.Server with profiling.
	pc := &profsvc.ProfilingConfig{
		Server:          s,
		Enabled:         true,
		StreamStatsSize: 1024,
	}
	if err = profsvc.Init(pc); err != nil {
		fmt.Printf("error calling profsvc.Init: %v\n", err)
		return err
	}

	go s.Serve(lis)
	return nil
}

func main() {
	flag.Parse()

	if err := setupClientProfiling(); err != nil {
		log.Fatalf("error setting up profiling: %v\n", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := echo.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.Simple(ctx, &echo.Request{Message: "hello, profiling"})
	fmt.Printf("Simple call returned %q, %v\n", res.GetMessage(), err)
	if err != nil {
		log.Fatalf("error calling Simple: %v", err)
	}

	log.Printf("sleeping for 30 seconds with exposed profiling service on :%d\n", *profilingPort)
	time.Sleep(30 * time.Second)
}
