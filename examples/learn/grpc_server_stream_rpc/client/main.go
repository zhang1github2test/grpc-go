package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/learn/grpc_server_stream_rpc/api"
	"google.golang.org/grpc/resolver"
	"io"
	"log"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//}
	resolver.SetDefaultScheme("passthrough")
	conn, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewWeatherServiceClient(conn)

	stream, err := client.ListWeather(context.Background(), &pb.WeatherRequest{
		Day: "2024-06-26",
	})
	if err != nil {
		log.Fatalf("failed to got list weather! ")
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("消息已经发送完成!")
			break
		}
		fmt.Println("响应回来的信息为:", resp)
	}

}
