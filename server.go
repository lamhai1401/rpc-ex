package main

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

func runServer() {
	rpc.RegisterName("KVStoreService", NewKVStoreService())
	rpc.RegisterName("HelloService", new(HelloService))
	// chạy rpc server trên port 1234
	listener, err := net.Listen("tcp", ":1234")
	// nếu có lỗi thì in ra
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	// vòng lặp để phục vụ nhiều client
	for {
		// accept một connection đến
		conn, err := listener.Accept()
		// in ra lỗi nếu có
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		// phục vụ client trên một goroutine khác
		// để giải phóng main thread tiếp tục vòng lặp
		rpc.ServeConn(conn)
		conn.Close()
	}
}

func runServeCallLient() {
	rpc.Register(new(HelloService))

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
