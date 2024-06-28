package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	"log"
	"time"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

const fallbackToken = "some-secret-token"

func callUnaryEcho(client ecpb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &ecpb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func unaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	log.Printf("rpc start %s, start time:%s", method, start.Format(time.RFC3339))
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	log.Printf("rpc end %s, end time:%s", method, end.Format(time.RFC3339))
	return err
}

func main() {
	flag.Parse()

	// 创建基于tls的凭证
	creds, err := credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.bac.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// 设置连接到服务的一个conn
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(creds), grpc.WithUnaryInterceptor(unaryInterceptor))

	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	defer conn.Close()
	// 创建一个echo client
	client := ecpb.NewEchoClient(conn)
	// 发送hello world消息到
	callUnaryEcho(client, "hello world!")
}
