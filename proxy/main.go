package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lamhai1401/rpc-ex/security"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

// func RunGateWay() {
// 	// khai báo biến context để xử lý signal kết thúc goroutine
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	// hàm cancel() sẽ kích hoạt ctx.Done()
// 	defer cancel()
// 	// mux được dùng cho việc routing
// 	mux := runtime.NewServeMux()

// 	// gọi hàm để đăng kí RestService cho proxy
// 	err := security.RegisterRestServiceHandlerFromEndpoint(
// 		// truyền vào biến ctx, mux, và địa chỉ gRPC service
// 		ctx, mux, "localhost:5000",
// 		[]grpc.DialOption{grpc.WithInsecure()},
// 	)

// 	// in ra lỗi nếu có
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// bắt đầu lắng nghe http client trên port 8080
// 	listener, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", zap.Error(err))
// 	}
// 	err = http.Serve(listener, h2c.NewHandler(
// 		mux,
// 		&http2.Server{}),
// 	)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", zap.Error(err))
// 	}
// }

func HttpGrpcRouter(httpHandler *runtime.ServeMux, listener net.Listener) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpHandler.ServeHTTP(w, r)
	})
}

func MakeHttpServer(conn *grpc.ClientConn) *runtime.ServeMux {
	router := runtime.NewServeMux()

	err := security.RegisterRestServiceHandlerFromEndpoint(
		// truyền vào biến ctx, mux, và địa chỉ gRPC service
		context.Background(), router, "localhost:5000",
		[]grpc.DialOption{grpc.WithInsecure()},
	)

	// err := security.RegisterRestServiceHandler(context.Background(), router, conn)
	if err != nil {
		panic(err.Error())
	}

	return router
}

func Start() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", zap.Error(err))
	}
	// serviced.MakeGrpcServer()
	router := MakeHttpServer(nil)

	fmt.Println("Starting server on address : " + ":8080")
	// err = http.Serve(listener, HttpGrpcRouter(router, listener))
	err = http.Serve(listener, h2c.NewHandler(
		HttpGrpcRouter(router, listener),
		&http2.Server{}),
	)
	return err
}
