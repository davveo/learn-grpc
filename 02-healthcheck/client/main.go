package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/davveo/learn-grpc/proto/echo"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

var serviceConfig = `{
	"loadBalancingPolicy": "round_robin",
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

func callSimple(c echo.EchoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Simple(ctx, &echo.Request{})
	if err != nil {
		fmt.Println("Simple: _, ", err)
	} else {
		fmt.Println("Simple: ", r.GetMessage())
	}
}

func main() {
	flag.Parse()

	r := manual.NewBuilderWithScheme("whatever")
	r.InitialState(resolver.State{
		Addresses: []resolver.Address{
			{Addr: "localhost:50051"},
			{Addr: "localhost:50052"},
		},
	})

	address := fmt.Sprintf("%s:///unused", r.Scheme())

	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithResolvers(r),
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	conn, err := grpc.Dial(address, options...)
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()

	echoClient := echo.NewEchoClient(conn)

	for {
		callSimple(echoClient)
		time.Sleep(time.Second)
	}
}
