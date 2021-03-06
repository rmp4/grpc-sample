package main

import (
	"context"
	"grpc-sample/pkg/gtls"
	"grpc-sample/pkg/pb"
	"log"
	"time"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const PORT = "9001"

func main() {
	tlsClient := gtls.Client{
		ServerName: "*.testserver.com",
		CaFile:     "./configs/ca-cert.pem",
		CertFile:   "./configs/client/client-cert.pem",
		KeyFile:    "./configs/client/client-key.pem",
	}

	c, err := tlsClient.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}
	auth := gtls.Authentication{
		User:     "admin",
		Password: "123",
	}

	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(5*time.Second)))
	defer cancel()

	client := pb.NewSearchServiceClient(conn)

	resp, err := client.Search(ctx, &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.Search err: deadline")
			}
		}

		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}
