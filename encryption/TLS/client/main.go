package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/davveo/learn-grpc/proto/echo"
	"google.golang.org/grpc"

	"github.com/davveo/learn-grpc/testdata"
	"google.golang.org/grpc/credentials"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func callSimple(client echo.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		resp, err := client.Simple(ctx, &echo.Request{Message: message})
		if err != nil {
			log.Fatalf("client.Simple(_) = _, %v: ", err)
		}
		fmt.Println("Simple: ", resp.Message)
	}

}

func main() {
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile(
		testdata.Path("ca.pem"), "x.test.youtube.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := echo.NewEchoClient(conn)
	callSimple(client, "hello world")
}
