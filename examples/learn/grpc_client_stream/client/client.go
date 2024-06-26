package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/learn/grpc_client_stream_rpc/api"
	"google.golang.org/grpc/resolver"
	"log"
	"math/rand"
	"sync"
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
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go RunListWeather(client, wg)
	}
	wg.Wait()

}

func RunListWeather(client pb.WeatherServiceClient, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second)
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
	}
	// 消息发送完成,发送关闭指令
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("stream.CloseAndRecv failed: %v", err)
	}
	log.Printf("服务端返回的结果为:%v", reply)
}

func randomRequest() *pb.WeatherRequest {
	return &pb.WeatherRequest{
		RequestId:   rand.Int31(),
		Time:        time.Now().Format(time.DateTime),
		Temperature: rand.Float32()*(38-20) + 20,
	}

}
