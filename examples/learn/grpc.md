​	gRPC 是一个高性能、开源的 RPC 框架，它最初由 Google 开发并开源。它使用 HTTP/2 作为传输协议，支持多种编程语言，并且通过 Protocol Buffers（protobuf）来定义接口和数据结构。以下是快速入门 gRPC 的步骤：

### 1、环境安装

[grpc快速入门][https://grpc.io/docs/languages/go/quickstart/]

1. 安装 gRPC 和 Protocol Buffers
   首先，你需要安装 Protocol Buffers 编译器 protoc 和 gRPC 库。

* 安装 Protocol Buffers 编译器
  从[官方页面][https://github.com/protocolbuffers/protobuf/releases]下载适合你操作系统的 protoc 编译器。

  ![image-20240624173650744](grpc.assets/image-20240624173650744.png)

下载好后解压，然后protoc的对应路径配置到环境变量中：D:\Users\admin\Downloads\protoc-gen-go.v1.34.2.windows.amd64

安装：

```go
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```



* 安装 gRPC 库
  安装 gRPC 库和 Protocol Buffers 插件

  ```
  go get -u google.golang.org/grpc
  go get -u github.com/golang/protobuf/protoc-gen-go
  ```

### 2、定义 gRPC 服务

创建一个 `.proto` 文件来定义你的 gRPC 服务和消息类型。以下是一个示例 calculator.proto 文件：

```proto
syntax = "proto3";

option go_package = "/api";

// The Calculator service definition.
service Calculator {
  // Sends a greeting
  rpc Add (AddRequest) returns (AddReply) {}
}

// The request message containing the user's name.
message AddRequest {
  int32 a = 1;
  int32 b = 2;
}

// The response message containing the Calculator
message AddReply {
  int32 reply = 1;
}

```

### 3、编译 calculator.proto 文件

```shell
protoc --go_out=. --go-grpc_out=. calculator.proto
```

### 4、编写客户端代码

```go
package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	pb "gorm-demo/grpc_calculator/api"
	"log"
	"time"
)

var (
	address = "localhost:50051"
)

func main() {
	// 不写这行会导致访问超时
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Add(ctx, &pb.AddRequest{
		A: 10,
		B: 20,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", r.GetReply())
}

```

### 5、编写服务端代码



核心代码为：

	lis, err := net.Listen("tcp", port)
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	s.Serve(lis);
```go
package main

import (
	"context"
	"google.golang.org/grpc"
	pb "gorm-demo/grpc_calculator/api"
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

```

