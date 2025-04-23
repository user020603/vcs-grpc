# Báo cáo GRPC

## 1. Khái niệm GRPC

- **Khái niệm**: gRPC (Google Remote Procedure Call) là một framework RPC hiện đại, mã nguồn mở, hiệu năng cao do Google phát triển. Nó cho phép ứng dụng client gọi trực tiếp các phương thức trên một ứng dụng server đặt ở máy khác như thể đó là một đối tượng cục bộ.

- **Nền tảng**: Hoạt động dựa trên HTTP/2, tận dụng các ưu điểm như multiplexing(cho phép nhiều luồng dữ liệu(requests/responses) được gửi đồng thời trên cùng một kết nối TCP), truyền dữ liệu nhị phân, nén header để tăng hiệu quả và giảm độ trễ.

- **Mục tiêu**: Được thiết kế để kết nối các dịch vụ trong môi trường microservices và giữa các trung tâm dữ liệu một cách hiệu quả, đáng tin cậy và đa nền tảng.

- **Nguyên tắc hoạt động cơ bản**:

  - **Client-Server Model**: Hoạt động theo mô hình client-server. Client gửi yêu cầu đến server, server xử lý và trả về phản hồi.

  - **Contract-First**: Sử dụng Protocol Buffers(Protobuf) để định nghĩa cấu trúc dữ liệu (messages) và các phương thức dịch vụ (services) trong một file `.proto`. File này đóng vai trò như một hợp đồng giữa client và server.

  - **Code Generation**: Từ file `.proto`, gRPC toolchain có thể tự động sinh ra mã nguồn (stub/skeleton) cho cả client và server ở nhiều ngôn ngữ lập trình khác nhau.

  - **Transport Protocol**: Thường sử dụng **HTTP/2** làm giao thức truyền tải, mang lại nhiều lợi ích về hiệu năng so với HTTP/1.1 (được REST sử dụng phổ biến).

## 2. So sánh gRPC và REST

| Tính năng                            | gRPC                                                 | REST (thường dùng JSON qua HTTP/1.1)                  |
| :----------------------------------- | :--------------------------------------------------- | :---------------------------------------------------- |
| **Định dạng dữ liệu**                | Protocol Buffers (nhị phân, nhỏ gọn)                 | JSON (text, dễ đọc bởi người)                         |
| **Giao thức truyền tải**             | HTTP/2 (hiệu năng cao, streaming)                    | Thường là HTTP/1.1 (đơn giản hơn)                     |
| **Định nghĩa API**                   | Contract-first (`.proto` file, chặt chẽ)             | Thường là Code-first (linh hoạt hơn)                  |
| **Kiểu giao tiếp**                   | Hỗ trợ Unary, Client/Server/Bi-directional Streaming | Thường là Request-Response (Unary)                    |
| **Code Generation**                  | Mạnh mẽ, đa ngôn ngữ                                 | Phụ thuộc vào các công cụ bên ngoài (Swagger/OpenAPI) |
| **Hiệu năng**                        | Cao hơn (nhờ Protobuf và HTTP/2)                     | Thấp hơn (do JSON và HTTP/1.1)                        |
| **Khả năng tương thích trình duyệt** | Hạn chế (cần gRPC-Web proxy)                         | Rất tốt                                               |
| **Độ phức tạp**                      | Cao hơn một chút ban đầu                             | Thấp hơn                                              |

### Khi nào chọn?

- **Chọn gRPC khi**:

  - Cần hiệu năng cao, độ trễ thấp (giao tiếp giữa các microservice).
  - Cần các kiểu streaming phức tạp.
  - Môi trường đa ngôn ngữ.
  - API nội bộ, không cần public trực tiếp ra trình duyệt.

- **Chọn REST khi**:
  - Cần API public, dễ dàng truy cập từ trình duyệt hoặc các client đơn giản.
  - Ưu tiên sự đơn giản và dễ đọc của dữ liệu (JSON).
  - Hệ sinh thái công cụ và thư viện hỗ trợ là ưu tiên hàng đầu.

## 3. Protocol Buffers (Protobuf)

Protocol Buffers là một cơ chế tuần tự hóa (serialization) dữ liệu có cấu trúc, được phát triển bởi Google. Nó trung lập về ngôn ngữ và nền tảng, có thể mở rộng, và được thiết kế để nhỏ gọn và nhanh chóng. Protobuf thường được sử dụng làm định dạng giao diện (Interface Definition Language - IDL) cho gRPC.

### Định nghĩa Message và Service

Định nghĩa cấu trúc dữ liệu (`message`) và các phương thức dịch vụ (`service`) trong một file có đuôi `.proto`.

```protobuf
syntax = "proto3";

package greet;

service Greeter {
    // Unary
    rpc SayHello(HelloRequest) returns (HelloReply);
}

message HelloRequest {
    string name = 1;
}

message HelloReply {
    string message = 1;
}
```

### Serialize/Deserialize Process

1. **Serialize (Client)**: Khi client gọi một phương thức stub RPC, dữ liệu yêu cầu (`HelloRequest`) được serialize thành định dạng nhị phân của Protobuf.

2. **Truyền tải**: Dữ liệu nhị phân được gửi qua mạng (thường qua HTTP/2).

3. **Deserialize (Server)**: Server nhận dữ liệu nhị phân và deserialize nó trở lại thành đối tượng `HelloRequest` mà code server skeleton có thể hiểu và xử lý.

4. **Serialize (Server)**: Sau khi hoàn tất xử lý, server serialize dữ liệu phản hồi (`HelloReply`) thành định dạng nhị phân.

5. **Truyền tải**: Dữ liệu nhị phân được gửi trả về Client.

6. **Deserialize (Client)**: Client nhận dữ liệu nhị phân, tiến hành deserialize nó thành đối tượng `HelloReply`.

## 4. Các loại RPC trong gRPC

gRPC định nghĩa 4 kiểu phương thức dịch vụ, dựa trên việc client và server gửi một hay nhiều message:

- **Unary RPC**:

  - Client gửi một yêu cầu duy nhất đến server.
  - Server xử lý và trả về một phản hồi duy nhất.
  - Giống mô hình request-response truyền thống của REST.
  - ex: `rpc SayHello(HelloRequest) returns (HelloReply);`

- **Server streaming RPC**:

  - Clienet gửi một yêu cầu duy nhất đến server.
  - Server xử lý và trả về một luồng (stream) các phản hồi. Client đọc từ luồng này cho đến khi khong còn message nào.
  - ex: `rpc ListFeatures(Rectangle) returns (stream Feature);` (Client gửi 1 khu vực, server trả về danh sách các địa điểm trong khu vực đó).

- **Client streaming RPC**:

  - Client gửi một luồng các yêu cầu đến server. Server đọc từ luồng này.
  - Khi client gửi xong, server xử lý và trả về một phản hồi tổng hợp duy nhất.
  - ex: `rpc RecordRoute(stream Point) returns (RouteSummary);` (Client gửi liên tục vị trí, server tính toán và trả về tóm tắt lộ trình khi client gửi xong).

- **Bidirectional streaming RPC**:
  - Cả client và server đều gửi một luồng các message cho nhau.
  - Hai luồng hoạt động độc lập, client và server có thể đọc/ghi theo bất kỳ thứ tự nào.
  - Phù hợp cho các ứng dụng tương tác cao như chat.
  - ex: `rpc RouteChat(stream RouteNote) returns (stream RouteNote);`

## 5. Xử lý lỗi (Error Handling)

Trong Go, gRPC xử lý lỗi thông qua package `google.golang.org/grpc/status` và các mã lỗi chuẩn codes.

```go
import (
    "context"
    "fmt"
    "google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
)

func callGetUser(client pb.UserServiceClient) {
    res, err := client.GetUser(context.Background(), &pb.UserRequest{Id: -1})
    if err != nil {
        st, ok := status.FromError(err)
        if ok {
            fmt.Println("❌ gRPC Error:")
            fmt.Println("Code:", st.Code())         // INVALID_ARGUMENT
            fmt.Println("Message:", st.Message())   // "user_id must be greater than 0"
        } else {
            fmt.Println("Không phải lỗi gRPC:", err)
        }
        return
    }

    fmt.Println("✅ User:", res.Name)
}
```

* **Thành công**: `OK` (mã 0)
* **Lỗi**: Các mã khác `OK` biểu thị lỗi. Một số mã phổ biến: 
    * `CANCELLED`: Client hủy yêu cầu.
    * `UNKNOWN`: Lỗi không xác định từ phía server.
    * `INVALID_ARGUMENT`: Client cung cấp tham số không hợp lệ.
    * `DEADLINE_EXCEEDED`: Yêu cầu hết hạn trước khi hoàn thành.
    * `NOT_FOUND`: Không tìm thấy tài nguyên yêu cầu.
    * `PERMISSION_DENIED`: Client không có quyền thực hiện.
    * `UNAUTHENTICATED`: Cần xác thực.
    * `UNAVAILABLE`: Dịch vụ tạm thời không khả dụng.
* **Metadata lỗi**: Ngoài mã trạng thái, server có thể gửi thêm thông tin chi tiết về lỗi dưới dạng metadata (key-value pairs). Client có thể đọc metadata này để hiểu rõ hơn về nguyên nhân lỗi. 

## 6. Bảo mật với SSL/TLS

gRPC tích hợp chặt chẽ với SSL/TLS để mã hóa dữ liệu truyền giữa client và server, đảm bảo tính bí mật và toàn vẹn. 

### Cơ chế hoạt động

1. Server Authentication 
- Server cung cấp chứng chỉ SSL/TLS để xác minh danh tính.
- Máy khách xác mình chứng chỉ này với một CA (Certificate Authority).

2. Thiết lập kết nối mã hóa
- Sau khi xác thực, client và server thực hiện bắt tay TLS (TLS handshake).
- Quá trình này tạo ra các khóa phiên được sử dụng để mã hóa dữ liệu. 

3. Truyền dữ liệu an toàn
- Toàn bộ giao tiếp giữa client và server được mã hóa.
- Dữ liệu được bảo vệ khỏi nghe lén và tấn công trung gian (man-in-the-middle).

### Các loại xác thực trong gRPC

1. Xác thực một chiều (One-way authentication)
- Client xác thực server thông qua chứng chỉ.
- Phổ biến trong nhiều ứng dụng web và API. 

2. Xác thực hai chiều (Two-way/mTLS)
- Cả client/server đều xác thực lẫn nhau.
- Client cũng cung cấp chứng chỉ cho server xác thực.
- Cung cấp mức độ bảo mật cao hơn cho các hệ thống nhạy cảm. 

## 7. gRPC Interceptor

Interceptor là một cơ chế cho phép chặn và xử lý các yêu cầu RPC đến (trên server) hoặc đi (trên client) trước khi chúng được xử lý thực sự hoặc gửi đi, hoạt động tương tự middleware trong ứng dụng REST API. 

* **Ứng dụng**:
    * **Logging**: Ghi log chi tiết về các request/response.
    * **Authentication/Authorization**: Kiểm tra thông tin các thực, phân quyền trước khi xử lý các logic nghiệp vụ.
    * **Monitoring/Metric**: Thu thập số liệu về thời gian xử lý, tỷ lệ xảy ra lỗi.
    * **Chain Interceptor**: Có thể kết hợp nhiều interceptor lại với nhau.

```go
func loggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    log.Printf("Received request: %s", info.FullMethod)
    start := time.Now()
    
    resp, err := handler(ctx, req)
    
    duration := time.Since(start)
    log.Printf("Request %s completed in %v", info.FullMethod, duration)
    
    return resp, err
}

func authUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Errorf(codes.Unauthenticated, "metadata not existed")
    }
    
    authHeader, ok := md["authorization"]
    if !ok || len(authHeader) == 0 {
        return nil, status.Errorf(codes.Unauthenticated, "token not existed")
    }
    
    token := authHeader[0]
    if !isValidToken(token) {
        return nil, status.Errorf(codes.Unauthenticated, "token not valided")
    }
    
    return handler(ctx, req)
}

server := grpc.NewServer(
    grpc.UnaryInterceptor(ChainUnaryInterceptors(
        loggingUnaryInterceptor,
        authUnaryInterceptor,
    )),
)
```

## 8. gRPC Reflection 

- gRPC Reflection là một dịch vụ tùy chọn mà server có thể kích hoạt. Khi được bật, nó cho phép các client có thể khám phá các service và method và message mà server cung cấp - mà không cần file `.proto`phía client. 
- Một trong những công cụ phổ biến nhất hỗ trợ gRPC Reflection là `Evans CLI`– một công cụ dòng lệnh có khả năng hoạt động tương tự như `curl` (dành cho HTTP) hoặc **Postman** (với gRPC). Với Evans, developer có thể duyệt service, gọi thử các RPC, và kiểm tra phản hồi một cách nhanh chóng và thuận tiện.

```go
import "google.golang.org/grpc/reflection"

func main() {
    server := grpc.NewServer()
    pb.RegisterYourServiceServer(server, &yourService{})
    reflection.Register(server)
    ...
}
```

```shell
evans --host localhost --port 50051 -r repl
```

## 9. Health Checking và Service Discovery

Trong môi trường microservices, việc biết được một service có đang hoạt động bình thường hay không và làm thế nào để tìm ra địa chỉ của nó là rất quan trọng.

*   **Health Checking:**
    *   gRPC định nghĩa một **Health Checking Protocol** chuẩn (`grpc.health.v1.Health`).
    *   Các service có thể implement protocol này để cung cấp endpoint (`Check`, `Watch`) cho phép các hệ thống khác (load balancer, orchestrator như Kubernetes) kiểm tra tình trạng sức khỏe của chúng (SERVING, NOT_SERVING, UNKNOWN).
*   **Service Discovery:**
    *   Là cơ chế để client tìm ra địa chỉ mạng (IP, port) của các instance server gRPC.
    *   gRPC không tự cung cấp một giải pháp service discovery cụ thể, nhưng nó tích hợp tốt với các hệ thống phổ biến như:
        *   **DNS:** Cách đơn giản nhất.
        *   **Load Balancer:** Client kết nối đến load balancer, load balancer điều phối đến các server instance (thường kết hợp health checking).
        *   **Service Registry:** Các công cụ như Consul, etcd, Zookeeper. Server đăng ký địa chỉ của mình vào registry, client truy vấn registry để tìm server.
