package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/learn/grpc_bidirectional_rpc/api"
	"google.golang.org/grpc/resolver"
	"log"
	"math/rand"
	"time"
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
	RunListWeather(client)

}

func RunListWeather(client pb.WeatherServiceClient) {
	var requests []*pb.WeatherRequest
	for i := 0; i < 10; i++ {
		requests = append(requests, randomRequest())
	}

	stream, err := client.ListWeather(context.Background())
	if err != nil {
		log.Fatalf("client.ListWeather failed, %v", err)
	}
	// 发送信息到服务端
	for _, req := range requests {
		stream.Send(req)
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		time.Sleep(time.Second)
		log.Println("获取到结果为", resp)
	}
	// 消息发送完成,发送关闭指令
	err = stream.CloseSend()
	if err != nil {
		log.Println(err)
	}
}

func randomRequest() *pb.WeatherRequest {
	return &pb.WeatherRequest{
		RequestId:   rand.Int31(),
		Time:        time.Now().Format(time.DateTime),
		Temperature: rand.Float32()*(38-20) + 20,
	}

}
