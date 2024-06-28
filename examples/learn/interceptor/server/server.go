package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"google.golang.org/grpc/examples/learn/echo/echo"
	"io"
	"log"
	"net"
	"time"
)

var port = flag.Int("port", 50051, "the port sever on")

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func (s *EchoServer) UnaryEcho(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	log.Printf("reviced request, msg is :%v", in)
	return &echo.EchoResponse{
		Message: in.Message,
	}, nil
}

func (s *EchoServer) ServerStreamingEcho(in *echo.EchoRequest, stream echo.Echo_ServerStreamingEchoServer) error {
	log.Printf("--ServerStreamingEcho--, request is:%v", in)
	for i := 0; i < 10; i++ {
		stream.Send(&echo.EchoResponse{
			Message: in.Message,
		})
	}
	return nil
}

func (s *EchoServer) ClientStreamingEcho(stream echo.Echo_ClientStreamingEchoServer) error {
	var rpcStatus error
	var message string
	for {
		recv, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			if rpcStatus == io.EOF {
				fmt.Printf("echo last received message\n")
				return stream.SendAndClose(&echo.EchoResponse{Message: message})
			}
		}
		message = recv.Message
		log.Printf("revice message from client: %v", recv)
	}
}

func (s *EchoServer) BidirectionalStreamingEcho(stream echo.Echo_BidirectionalStreamingEchoServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending echo\n", in)
		if err := stream.Send(&echo.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	startTime := time.Now()
	log.Printf("start to proccess, method:%v, time is:%v", info.FullMethod, startTime.Format(time.RFC3339Nano))
	resp, err = handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed with error: %v", err)
	}
	endTime := time.Now()
	log.Printf("end to proccess, method:%v, startTime is:%v, endTime is:%v, Spending time:%v", info.FullMethod,
		startTime.Format(time.RFC3339Nano), endTime.Format(time.RFC3339Nano), endTime.Sub(startTime))
	return
}

func main() {
	// 加载TLS文件
	creds, err := credentials.NewServerTLSFromFile(data.Path("x509/server_cert.pem"), data.Path("x509/server_key.pem"))
	if err != nil {
		log.Fatalf("load creds failed:%v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds), grpc.ChainUnaryInterceptor(UnaryServerInterceptor))
	echo.RegisterEchoServer(s, &EchoServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
