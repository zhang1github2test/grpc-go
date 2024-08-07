package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"google.golang.org/grpc/examples/learn/echo/echo"
	"google.golang.org/grpc/resolver"
	"io"
	"log"
	"time"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	resolver.SetDefaultScheme("passthrough")
	// 加载TLS证书
	creds, err := credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load creds: %v", err)
	}

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("failed to conn to server: %v", err)
	}
	defer conn.Close()
	client := echo.NewEchoClient(conn)

	// 一元方法客户端使用
	callUnaryEcho(client, "hello world!")

	// 调用服务端流式方法
	callServerStream(client, "helo, world!")

	// 调用客户端流式方法
	callClientStream(client, "hello, world!")

	//调用双向流式方法
	callBidirectionalStreamingEcho(client, "hello world!")

}

func callUnaryEcho(client echo.EchoClient, msg string) {
	log.Printf("--- UnaryEcho start---\n")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()
	unaryEcho, err := client.UnaryEcho(ctx, &echo.EchoRequest{Message: msg})

	if err != nil {
		log.Fatalf("client.UnaryEcho(_)=%v", err)
	}
	log.Printf("--- UnaryEcho end---UnaryEcho response:%v\n", unaryEcho)
}

func callServerStream(client echo.EchoClient, msg string) {
	log.Printf("--- server streaming start---\n")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	stream, err := client.ServerStreamingEcho(ctx, &echo.EchoRequest{
		Message: msg,
	})
	if err != nil {
		log.Fatalf("failed to call ServerStreamingEcho: %v", err)
	}
	var rpcStatus error
	for {
		recv, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		log.Printf("response:%v", recv)
	}
	if rpcStatus != io.EOF && rpcStatus != nil {
		log.Fatalf("failed to finish server streaming: %v", rpcStatus)
	}
	log.Println("--- server streaming end---")

}

func callClientStream(client echo.EchoClient, msg string) {
	log.Printf("--- Client Streaming Echo start---\n")
	stream, err := client.ClientStreamingEcho(context.Background())
	if err != nil {
		log.Fatalf("failed to call ClientStreamingEcho!")
	}
	// 发送多个消息给服务端
	for i := 0; i < 10; i++ {
		err := stream.Send(&echo.EchoRequest{
			Message: msg,
		})
		if err != nil {
			log.Fatalf("failed to send echo msg!")
		}
	}

	// 获取服务端发回来的信息
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to call CloseAndRecv!%v", err)
	}
	log.Printf("--- Client Streaming Echo end--- reponse:%v", resp)

}

func callBidirectionalStreamingEcho(client echo.EchoClient, msg string) {
	log.Printf("--- Bidirectional Streaming Echo ---\n")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream, err := client.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("falied to call BidirectionalStreamingEcho: %v\n", err)
	}

	// 往客户发送流式消息
	go func() {
		for i := 0; i < 10; i++ {
			if err := stream.Send(&echo.EchoRequest{Message: msg}); err != nil {
				log.Fatalf("failed to send msg: %v\n", err)
			}
		}
		// 发送完成后关闭stream
		stream.CloseSend()
	}()
	// 读取从服务端返回的流式信息
	var rpcStatus error
	for i := 0; i < 10; i++ {
		resp, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		if err != nil {
			log.Fatalf("failed to recv msg: %v\n", err)
		}
		log.Printf("reviced msg is:%v\n", resp)
	}
	if rpcStatus != io.EOF && rpcStatus != nil {
		log.Fatalf("failed to finish server streaming: %v", rpcStatus)
	}
	log.Printf("--- Bidirectional Streaming Echo end---")

}
