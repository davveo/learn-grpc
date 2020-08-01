package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/davveo/learn-grpc/proto/echo"

	"github.com/davveo/learn-grpc/testdata"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}

func callSimple(client echo.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Simple(ctx, &echo.Request{Message: message})
	if err != nil {
		log.Fatalf("client.callSimple(_) = _, %v: ", err)
	}
	fmt.Println("callSimple: ", resp.Message)
}

func main() {
	flag.Parse()
	perRPC := oauth.NewOauthAccess(fetchToken())
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("ca.pem"), "x.test.youtube.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(creds),
	}
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rgc := echo.NewEchoClient(conn)

	callSimple(rgc, "hello world")
}
