package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/learn/grpc_server_stream_rpc/api"
	"log"
	"math/rand"
	"net"
	"time"
)

type WeatherServiceServer struct {
	pb.UnimplementedWeatherServiceServer
}

func (WeatherServiceServer) ListWeather(request *pb.WeatherRequest, stream pb.WeatherService_ListWeatherServer) error {
	rand.NewSource(time.Now().UnixNano())
	day := request.GetDay()
	startTime, err := time.Parse("2006-01-02", day)
	if err != nil {
		stream.SendMsg("传入的参不正确" + err.Error())
	}
	endTime := startTime.AddDate(0, 0, 1)

	for ; startTime.Before(endTime); startTime = startTime.Add(time.Hour) {
		resp := &pb.WeatherResponse{
			Temperature: rand.Float32()*(35-20) + 20,
			Timestamp:   startTime.Unix(),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
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
