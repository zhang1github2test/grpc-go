package main

import (
	"context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/learn/grpc_unary_rpc/api"
	"log"
	"net"
)

const (
	port = ":50051"
)

// server is used to implement Add.CalculatorServer.
type server struct {
	pb.UnimplementedCalculatorServer
}

// Add implements add.CalculatorServer
func (s *server) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddReply, error) {
	ctx.Done()
	log.Println("接收到计算请求")
	return &pb.AddReply{Reply: in.A + in.B}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
