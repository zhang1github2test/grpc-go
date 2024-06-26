package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/learn/grpc_unary_rpc/api"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

var (
	address = "localhost:50051"
)

func main() {
	// 不添加这行会导致访问超时
	resolver.SetDefaultScheme("passthrough")
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalculatorClient(conn)

	for i := 0; i < 100; i++ {
		time.Sleep(2 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.Add(ctx, &pb.AddRequest{
			A: 10,
			B: 20,
		})
		if err != nil {
			log.Printf("could not greet: %v", err)
			continue
		}
		log.Printf("Greeting: %v", r.GetReply())
	}

}
