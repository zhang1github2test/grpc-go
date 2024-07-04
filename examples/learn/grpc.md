​	gRPC 是一个高性能、开源的 RPC 框架，它最初由 Google 开发并开源。它使用 HTTP/2 作为传输协议，支持多种编程语言，并且通过 Protocol Buffers（protobuf）来定义接口和数据结构。以下是快速入门 gRPC 的步骤：

# 一、快速入门

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

# 二、GRPC的四种支持的四种方法模式

## 1. 简单RPC（Unary RPC）

### 描述

客户端发送单个请求到服务器，并接收单个响应。这是最常见的RPC模式。

### 应用场景

- 简单的查询和响应操作，如获取用户信息、查询数据库记录等。

  

## 2. 服务器端流式RPC（Server Streaming RPC）

### 描述

客户端发送一个请求到服务器，服务器返回一个流来发送一系列的消息。客户端从流中读取消息，直到流结束。

### 应用场景

- 需要持续提供数据的场景，如实时日志更新、数据推送等。

### 示例

演示客户端发送一个日期，服务端把每个小时的温度数据按流的方式将数据返回

```proto
// .proto文件定义
syntax = "proto3";
option go_package = "/api";

// The Weather service definition.
service WeatherService {
  // Sends a weather quest
  rpc ListWeather(WeatherRequest) returns (stream WeatherResponse) {}
}

// 天气请求对象
message WeatherRequest {
  // 要查询天气的日期 格式2006-01-02
  string day = 1;
}

// The response message containing the timestamp and temperature
message WeatherResponse {
  float temperature = 1;
  // 天气对应的时间戳
  int64 timestamp = 2;

}

```

```go
// 服务端实现

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

```

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	pb "gorm-demo/grpc_calculator_server_stream/api"
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
```





## 3. 客户端流式RPC（Client Streaming RPC）

### 描述

客户端通过流发送一系列的请求到服务器，服务器在接收完所有的请求后返回一个响应。

### 应用场景

- 需要批量上传数据的场景，如上传日志文件、批量提交数据等。

### 示例

```proto
// .proto文件定义
service Greeter {
  rpc RecordRoute(stream Point) returns (RouteSummary);
}
```

```go
// 服务端实现
func (s *server) RecordRoute(stream pb.Greeter_RecordRouteServer) error {
  var pointCount int32
  for {
    point, err := stream.Recv()
    if err == io.EOF {
      return stream.SendAndClose(&pb.RouteSummary{PointCount: pointCount})
    }
    if err != nil {
      return err
    }
    pointCount++
  }
}
```

## 4. 双向流式RPC（Bidirectional Streaming RPC）

### 描述

客户端和服务器都能发送一系列的消息，并且双方可以在任意顺序下读取和写入。这提供了最大程度的灵活性。

### 应用场景

- 需要实时双向通信的场景，如聊天应用、视频会议等。

### 示例

```proto
// .proto文件定义
service Greeter {
  rpc Chat(stream Message) returns (stream Message);
}
```

```go
// 服务端实现
func (s *server) Chat(stream pb.Greeter_ChatServer) error {
  for {
    msg, err := stream.Recv()
    if err == io.EOF {
      return nil
    }
    if err != nil {
      return err
    }
    if err := stream.Send(msg); err != nil {
      return err
    }
  }
}
```

## 总结

| 方法类型             | 描述          | 应用场景           |
| -------------------- | ------------- | ------------------ |
| 简单RPC（Unary RPC） | 单请求-单响应 | 查询、简单交互     |
| 服务器端流式RPC      | 单请求-多响应 | 数据推送、实时更新 |
| 客户端流式RPC        | 多请求-单响应 | 批量上传           |
| 双向流式RPC          | 多请求-多响应 | 实时双向通信       |

# 三、GRPC的metata数据

在 gRPC 中，Metadata 是一种用于在客户端和服务器之间传递附加信息的机制。它分为两种：

1. **Header Metadata**：在 RPC 调用开始时传递，用于传递身份验证、初始化参数等信息。
2. **Trailer Metadata**：在 RPC 调用结束时传递，用于传递状态码、处理结果等信息。

好的，以下是 Header Metadata 和 Trailer Metadata 的相同点和不同点的总结表格：

| 特性               | Header Metadata                  | Trailer Metadata              |
| ------------------ | -------------------------------- | ----------------------------- |
| **相同点**         |                                  |                               |
| 传递机制           | gRPC 的元数据传递机制的一部分    | gRPC 的元数据传递机制的一部分 |
| 数据类型           | 键值对 (Key-Value Pairs)         | 键值对 (Key-Value Pairs)      |
| 传递目的           | 传递附加信息                     | 传递附加信息                  |
| **不同点**         |                                  |                               |
| 传递时机           | RPC 调用开始时                   | RPC 调用结束时                |
| 用途               | 身份验证、初始化参数、追踪信息等 | 状态码、处理结果、统计信息等  |
| 设置位置（客户端） | `metadata.NewOutgoingContext`    | 在 RPC 调用之后获取           |
| 设置位置（服务器） | 在处理请求开始时设置             | 在处理请求结束时设置          |

这个表格展示了 Header Metadata 和 Trailer Metadata 的相同点和不同点，帮助理解它们在 gRPC 中的不同用途和使用时机。



# 四、GRPC的拦截器

​	gRPC 提供了简单的 API 来实现和安装拦截器，可以按每个 ClientConn/Server 进行设置。拦截器作为应用程序和 gRPC 之间的一层，可以用来观察或控制 gRPC 的行为。拦截器可以用于日志记录、身份验证/授权、指标收集以及其他跨 RPC 共享的功能。

## 1.客户端拦截器

###  一元拦截器(`UnaryClientInterceptor`)  



### 流拦截器 (`StreamClientInterceptor`)

## 2.服务端拦截器

###  一元拦截器(`UnaryClientInterceptor`)  

### 流拦截器 (`StreamClientInterceptor`)

| **类别** | **拦截器类型**                        | **函数签名**                                                 | **说明**                                                     | **安装方法**                                                 |
| -------- | ------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 客户端   | 一元拦截器 (`UnaryClientInterceptor`) | `func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error` | 实现通常分为三个部分：预处理、调用 RPC 方法、后处理。预处理可以检查和修改 RPC 调用参数，调用 `invoker` 执行 RPC 方法，后处理处理返回的回复和错误。 | 使用 `Dial` 配置 [`WithUnaryInterceptor`](https://godoc.org/google.golang.org/grpc#WithUnaryInterceptor) |
| 客户端   | 流拦截器 (`StreamClientInterceptor`)  | `func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)` | 实现通常包括预处理和流操作拦截。预处理类似于一元拦截器，流操作拦截通过包装 `ClientStream` 并重载其方法实现。 | 使用 `Dial` 配置 [`WithStreamInterceptor`](https://godoc.org/google.golang.org/grpc#WithStreamInterceptor) |
| 服务器端 | 一元拦截器 (`UnaryServerInterceptor`) | `func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error)` | 参考客户端一元拦截器的实现和说明。                           | 使用 `NewServer` 配置 [`UnaryInterceptor`](https://godoc.org/google.golang.org/grpc#UnaryInterceptor) |
| 服务器端 | 流拦截器 (`StreamServerInterceptor`)  | `func(srv interface{}, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error` | 参考客户端流拦截器的实现和说明。                             | 使用 `NewServer` 配置 [`StreamInterceptor`](https://godoc.org/google.golang.org/grpc#StreamInterceptor) |

​																															gRPC 拦截器对比



# 五、GRPC的启用压缩

| 优点                                                         | 缺点                                                         |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| **减少网络带宽使用**                                         | **增加 CPU 开销**                                            |
| GZIP 压缩显著减小数据包大小，减少网络传输所需带宽。          | 压缩和解压缩数据需要额外的 CPU 资源，可能导致性能下降。      |
| **提高传输速度**                                             | **延迟增加**                                                 |
| 在带宽限制条件下，数据量减少可以提高传输速度，使通信更快速。 | 压缩和解压缩本身会引入额外的延迟，特别是在数据量较小且网络带宽较高的情况下。 |
| **降低传输成本**                                             | **配置复杂性**                                               |
| 特别是在云服务中，压缩数据可以降低传输量，从而减少成本。     | 启用压缩需要额外的配置和代码修改，增加了开发和维护的复杂性。 |
| **优化用户体验**                                             | **兼容性问题**                                               |
| 对用户来说，较快的响应时间和更少的数据消耗意味着更好的用户体验，特别是在移动网络环境下。 | 并不是所有客户端和服务器都支持 GZIP 压缩，可能引发兼容性问题。 |

## 客户端配置：

在调用客户端方法的时候添加上grpc.UseCompressor(gzip.Name)。如下所示：

```
res, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: msg}, grpc.UseCompressor(gzip.Name))
```

## 服务端配置：

在import的代码端中安装gzip来进行解压和压缩,如下所示：

```go
_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
```

## 完整客户端的示例为：

```go
// Binary client is an example client.
package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/resolver"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip" // Install the gzip compressor
	pb "google.golang.org/grpc/examples/features/proto/echo"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()
	resolver.SetDefaultScheme("passthrough")
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)

	// Send the RPC compressed.  If all RPCs on a client should be sent this
	// way, use the DialOption:
	// grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name))
	var msg = "Send the RPC compressed.  If all RPCs on a client should be sent this" + "way, use the DialOption:"

	log.Println("原始消息长度", len(msg))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: msg}, grpc.UseCompressor(gzip.Name))
	fmt.Printf("UnaryEcho call returned %q, %v\n", res.GetMessage(), err)
	if err != nil || res.GetMessage() != msg {
		log.Fatalf("Message=%q, err=%v; want Message=%q, err=<nil>", res.GetMessage(), err, msg)
	}

	res, err = c.UnaryEcho(ctx, &pb.EchoRequest{Message: msg})
	fmt.Printf("UnaryEcho call returned %q, %v\n", res.GetMessage(), err)
	if err != nil || res.GetMessage() != msg {
		log.Fatalf("Message=%q, err=%v; want Message=%q, err=<nil>", res.GetMessage(), err, msg)
	}

}

```

## 完整服务端示例为：

```go
// Binary server is an example server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor

	pb "google.golang.org/grpc/examples/features/proto/echo"
)

var port = flag.Int("port", 50051, "the port to serve on")

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("UnaryEcho called with message %q\n", in.GetMessage())
	return &pb.EchoResponse{Message: in.Message}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}

```

