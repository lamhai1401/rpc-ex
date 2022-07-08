package pubsub

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/lamhai1401/rpc-ex/model"
	"github.com/moby/moby/pkg/pubsub"
	"google.golang.org/grpc"
)

func Run() {
	// khởi tạo một đối tượng gRPC service
	grpcServer := grpc.NewServer()

	// đăng ký service với grpcServer (của gRPC plugin)
	model.RegisterPubsubServiceServer(grpcServer, NewPubsubService())

	// cung cấp gRPC service trên port `1234`
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}

type PubsubService struct {
	pub *pubsub.Publisher
}

func NewPubsubService() *PubsubService {
	return &PubsubService{
		pub: pubsub.NewPublisher(100*time.Millisecond, 10),
	}
}

func (p *PubsubService) Publish(
	ctx context.Context, arg *model.String,
) (*model.String, error) {
	p.pub.Publish(arg.GetValue())
	return &model.String{}, nil
}

func (p *PubsubService) Subscribe(
	arg *model.String, stream model.PubsubService_SubscribeServer,
) error {
	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, arg.GetValue()) {
				return true
			}
		}
		return false
	})

	for v := range ch {
		if err := stream.Send(&model.String{Value: v.(string)}); err != nil {
			return err
		}
	}

	return nil
}
