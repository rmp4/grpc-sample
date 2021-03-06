package main

import (
	"context"
	"grpc-sample/pkg/gtls"
	"grpc-sample/pkg/pb"
	"log"
	"net"
	"runtime/debug"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchService struct {
	auth *gtls.Authentication
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	auth := gtls.Authentication{
		User:     "admin",
		Password: "123",
	}
	s.auth = &auth
	// Need to add on every function
	if err := s.auth.Auth(ctx); err != nil {
		return nil, err
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "searchService.Search canceled")
	}

	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	tlsServer := gtls.Server{
		CaFile:   "./configs/ca-cert.pem",
		CertFile: "./configs/server/server-cert.pem",
		KeyFile:  "./configs/server/server-key.pem",
	}
	c, err := tlsServer.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}
