package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/davveo/learn-grpc/proto/echo"
	"github.com/davveo/learn-grpc/proto/helloworld"

	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func callSayHello(c helloworld.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("client.SayHello(_) = _, %v", err)
	}
	fmt.Println("Greeting: ", r.Message)
}

func callSimple(client echo.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.Simple(ctx, &echo.Request{Message: message})
	if err != nil {
		log.Fatalf("client.Simple(_) = _, %v: ", err)
	}
	fmt.Println("Simple: ", resp.Message)
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	fmt.Println("--- calling helloworld.Greeter/SayHello ---")
	// Make a greeter client and send an RPC.
	hwc := helloworld.NewGreeterClient(conn)
	callSayHello(hwc, "multiplex")

	fmt.Println()
	fmt.Println("--- calling routeguide.RouteGuide/GetFeature ---")
	// Make a routeguild client with the same ClientConn.
	rgc := echo.NewEchoClient(conn)
	callSimple(rgc, "this is examples/multiplex")
}
