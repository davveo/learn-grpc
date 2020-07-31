package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	echo "github.com/davveo/learn-grpc/pb"
	"github.com/davveo/learn-grpc/testdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	echo.UnimplementedEchoServer
}

func (s *server) Simple(ctx context.Context, in *echo.Request) (*echo.Response, error) {
	return &echo.Response{Message: in.Message}, nil

}
func (s *server) ClientStream(echo.Echo_ClientStreamServer) error {
	return nil

}
func (s *server) ServerStream(*echo.Request, echo.Echo_ServerStreamServer) error {
	return nil
}
func (s *server) DoubleStream(echo.Echo_DoubleStreamServer) error {
	return nil
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// 这里应该从缓存中获取
	clientSecret := "some-secret-token"
	return token == clientSecret
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// 从ctx中获取授权信息
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// 如果正常获取，则执行后续逻辑
	return handler(ctx, req)
}

func main() {
	flag.Parse()
	fmt.Printf("server starting on port %d...\n", *port)
	cert, err := tls.LoadX509KeyPair(testdata.Path("server1.pem"), testdata.Path("server1.key"))
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}
	s := grpc.NewServer(opts...)
	echo.RegisterEchoServer(s, &server{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
