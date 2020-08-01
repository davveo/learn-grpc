package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/credentials/alts"

	echo "github.com/davveo/learn-grpc/pb"
	"google.golang.org/grpc"
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
	creds := alts.NewClientCreds(alts.DefaultClientOptions())
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := echo.NewEchoClient(conn)
	callSimple(client, "hello world")
}
