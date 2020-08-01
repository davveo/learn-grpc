package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/davveo/learn-grpc/proto/echo"

	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
	// 更多配置参考https://github.com/grpc/grpc/blob/master/doc/service_config.md
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "grpc.examples.echo.Echo"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)

// 配置信息
func retryDial() (*grpc.ClientConn, error) {
	return grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(retryPolicy))
}

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := retryDial()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()

	c := echo.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for {
		reply, err := c.Simple(ctx, &echo.Request{Message: "Try and Success"})
		if err != nil {
			log.Fatalf("Simple error: %v", err)
		}
		log.Printf("Simple reply: %v", reply)
	}

}
