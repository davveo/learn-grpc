package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/davveo/learn-grpc/proto/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
)

var (
	ports = []string{":10001", ":10002", ":10003"}
	r     = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu    sync.Mutex
)

func Intn(n int) int {
	mu.Lock()
	res := r.Intn(n)
	mu.Unlock()
	return res
}

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

type slowServer struct {
	helloworld.UnimplementedGreeterServer
}

func (s *slowServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	// Delay 100ms ~ 200ms before replying
	time.Sleep(time.Duration(100+Intn(100)) * time.Millisecond)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	/***** Set up the server serving channelz service. *****/
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)
	go s.Serve(lis)
	defer s.Stop()

	/***** Start three GreeterServers(with one of them to be the slowServer). *****/
	for i := 0; i < 3; i++ {
		lis, err := net.Listen("tcp", ports[i])
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()
		s := grpc.NewServer()
		if i == 2 {
			helloworld.RegisterGreeterServer(s, &slowServer{})
		} else {
			helloworld.RegisterGreeterServer(s, &server{})
		}
		go s.Serve(lis)
	}

	/***** Wait for user exiting the program *****/
	select {}
}
