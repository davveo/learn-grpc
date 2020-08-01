package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/davveo/learn-grpc/proto/echo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var addr = flag.String("addr", "localhost:50052", "the address to connect to")

func unaryCall(c echo.EchoClient, requestID int, message string, want codes.Code) {
	// Creates a context with a one second deadline for the RPC.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &echo.Request{Message: message}

	_, err := c.Simple(ctx, req)
	got := status.Code(err)
	fmt.Printf("[%v] wanted = %v, got = %v\n", requestID, want, got)
}

func streamingCall(c echo.EchoClient, requestID int, message string, want codes.Code) {
	// Creates a context with a one second deadline for the RPC.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.DoubleStream(ctx)
	if err != nil {
		log.Printf("Stream err: %v", err)
		return
	}

	err = stream.Send(&echo.Request{Message: message})
	if err != nil {
		log.Printf("Send error: %v", err)
		return
	}

	_, err = stream.Recv()

	got := status.Code(err)
	fmt.Printf("[%v] wanted = %v, got = %v\n", requestID, want, got)
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := echo.NewEchoClient(conn)

	// A successful request
	unaryCall(c, 1, "world", codes.OK)
	// Exceeds deadline
	unaryCall(c, 2, "delay", codes.DeadlineExceeded)
	// A successful request with propagated deadline
	unaryCall(c, 3, "[propagate me]world", codes.OK)
	// Exceeds propagated deadline
	unaryCall(c, 4, "[propagate me][propagate me]world", codes.DeadlineExceeded)
	// Receives a response from the stream successfully.
	streamingCall(c, 5, "[propagate me]world", codes.OK)
	// Exceeds propagated deadline before receiving a response
	streamingCall(c, 6, "[propagate me][propagate me]world", codes.DeadlineExceeded)
}
