package pubsub

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/lamhai1401/rpc-ex/model"
	"google.golang.org/grpc"
)

func RunClientPub() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := model.NewPubsubServiceClient(conn)

	time.Sleep(1 * time.Second)

	_, err = client.Publish(
		context.Background(), &model.String{Value: "go: hello Go"},
	)

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	_, err = client.Publish(
		context.Background(), &model.String{Value: "docker: hello Docker"},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func RunClientSub() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := model.NewPubsubServiceClient(conn)

	stream1, err := client.Subscribe(context.Background(), &model.String{Value: "go:"})
	if err != nil {
		panic(err.Error())
	}

	stream2, _ := client.Subscribe(context.Background(), &model.String{Value: "docker:"})

	go func() {
		for {
			reply, err := stream1.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			fmt.Println("RunClientSub ", reply.GetValue())
		}
	}()

	go func() {
		for {
			reply, err := stream2.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			fmt.Println("RunClientSub ", reply.GetValue())
		}
	}()

	select {}
}
