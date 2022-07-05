package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

func runClient() {
	// kết nối đến rpc server
	client, err := rpc.Dial("tcp", "localhost:1234")
	// in ra lỗi nếu có
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// biến chứa giá trị trả về sau lời gọi rpc
	// var reply string
	// // gọi rpc với tên service đã register
	// err = client.Call("HelloService.Hello", "World", &reply)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	doClientWork(client)

	// // in ra kết quả
	// fmt.Println("reply:", reply)

	select {}
}

func doClientWork(client *rpc.Client) {
	// khởi chạy một Goroutine riêng biệt để giám sát khóa thay đổi
	go func() {
		var keyChanged string
		// lời gọi `watch` synchronous sẽ block cho đến khi
		// có khóa thay đổi hoặc timeout
		for {
			err := client.Call("KVStoreService.Watch", 30, &keyChanged)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("watch:", keyChanged)
		}
	}()

	time.Sleep(time.Second * 3)

	err := client.Call(
		//  giá trị KV được thay đổi bằng phương thức `Set`
		"KVStoreService.Set", [2]string{"abc", "abc-value"},
		new(struct{}),
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	time.Sleep(time.Second * 3)

	//  set lại lần nữa để giá trị value của key 'abc' thay đổi
	err = client.Call(
		"KVStoreService.Set", [2]string{"abc", "another-value"},
		new(struct{}),
	)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
}

func runRPCClient() {
	// listen trên port 1234 chờ server gọi
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	clientChan := make(chan *rpc.Client)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("Accept error:", err)
			}

			// khi mỗi đường link được thiết lập, đối tượng
			// RPC client được khởi tạo dựa trên link đó và
			// gửi tới client channel
			clientChan <- rpc.NewClient(conn)
		}
	}()

	doRPCClient(clientChan)
}

func doRPCClient(clientChan chan *rpc.Client) {
	//  nhận vào đối tượng RPC client từ channel
	client := <-clientChan

	// đóng kết nối với client trước khi hàm exit
	defer client.Close()

	var reply string

	// thực hiện lời gọi rpc bình thường
	err := client.Call("HelloService.Hello", "Lam Hai", &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
