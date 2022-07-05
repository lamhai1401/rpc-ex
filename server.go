package main

import (
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/lamhai1401/rpc-ex/model"
	"google.golang.org/grpc"
)

func runGRPCServer() {
	// khởi tạo một đối tượng gRPC service
	grpcServer := grpc.NewServer()

	// đăng ký service với grpcServer (của gRPC plugin)
	model.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	// cung cấp gRPC service trên port `1234`
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}

func runServeCallLient() {
	rpc.Register(new(KVStoreService))

	for {
		// chủ động gọi tới client
		conn, _ := net.Dial("tcp", "localhost:1234")
		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		rpc.ServeConn(conn)
		conn.Close()
	}
}
