package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/davveo/learn-grpc/proto/echo"
	"github.com/davveo/learn-grpc/utils"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func serverStream(c echo.EchoClient, message string) {

	stream, err := c.ServerStream(context.Background(), &echo.Request{Message: message})
	if err != nil {
		log.Fatalf("failed to call ServerStreamingEcho: %v", err)
	}

	var rpcStatus error
	fmt.Printf("response:\n")
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		fmt.Printf(" - %s\n", r.Message)
	}
	if rpcStatus != io.EOF {
		log.Fatalf("failed to finish server streaming: %v", rpcStatus)
	}
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	utils.CheckError(err)
	defer conn.Close()

	rgc := echo.NewEchoClient(conn)
	serverStream(rgc, "hello world")
}
