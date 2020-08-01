package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/davveo/learn-grpc/proto/echo"
	"github.com/davveo/learn-grpc/utils"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func callSimpleEcho(client echo.EchoClient, message string) {
	ctx := context.Background()
	resp, err := client.Simple(ctx, &echo.Request{Message: message})
	utils.CheckError(err)

	fmt.Println("simple echo: ", resp.Message)
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	utils.CheckError(err)
	defer conn.Close()

	rgc := echo.NewEchoClient(conn)
	callSimpleEcho(rgc, "hello world")
}
