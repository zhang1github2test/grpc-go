package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/learn/grpc_client_stream/api"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type WeatherServiceServer struct {
	pb.UnimplementedWeatherServiceServer
}

func (*WeatherServiceServer) ListWeather(stream pb.WeatherService_ListWeatherServer) error {
	var sum int32 = 0
	var errResult error
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			err := stream.SendAndClose(&pb.WeatherResponse{
				Process: true,
			})
			errResult = err
			break
		}

		//
		if err != nil {
			log.Println("通道关闭...", err.Error())
			errResult = err
			break
		}
		atomic.AddInt32(&sum, 1)
		fmt.Printf("接收到客户上报过来的温度信息,待进行到后续的操作%v,次数:%v.\n", req, sum)
		// ... 执行一些其他的业务操作,如将数据保存到数据库,温度过高产生报警信息等!
	}
	fmt.Println("执行次数", sum)

	return errResult
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpcServer, &WeatherServiceServer{})
	grpcServer.Serve(lis)
}
