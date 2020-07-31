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
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	streamingCount  = 10
)

func doubleStream(c echo.EchoClient, message string) {
	fmt.Printf("--- doubleStream ---\n")
	stream, err := c.DoubleStream(context.Background())
	if err != nil {
		log.Fatalf("failed to call BidirectionalStreamingEcho: %v\n", err)
	}

	go func() {
		for i:=0; i<streamingCount; i++ {
			if err := stream.Send(&echo.Request{Message: message}); err != nil {
				log.Fatalf("failed to send streaming: %v\n", err)
			}
		}
		stream.CloseSend()
	}()

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

func main()  {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	utils.CheckError(err)
	defer conn.Close()

	rgc := echo.NewEchoClient(conn)
	doubleStream(rgc, "hello world")
}