package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/davveo/learn-grpc/proto/echo"
	"github.com/davveo/learn-grpc/utils"
	"google.golang.org/grpc"
)

var (
	addr           = flag.String("addr", "localhost:50051", "the address to connect to")
	streamingCount = 10
)

func clientStream(c echo.EchoClient, message string) {
	fmt.Printf("--- client streaming ---\n")
	stream, err := c.ClientStream(context.Background())
	if err != nil {
		log.Fatalf("failed to call ClientStreamingEcho: %v\n", err)
	}

	for i := 0; i < streamingCount; i++ {
		log.Printf("send message: %d", i)
		if err := stream.Send(&echo.Request{Message: message}); err != nil {
			log.Fatalf("failed to send streaming: %v\n", err)
		}
	}

	r, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to CloseAndRecv: %v\n", err)
	}
	fmt.Printf("response:\n")
	fmt.Printf(" - %s\n\n", r.Message)
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	utils.CheckError(err)
	defer conn.Close()

	rgc := echo.NewEchoClient(conn)
	clientStream(rgc, "hello world")
}
