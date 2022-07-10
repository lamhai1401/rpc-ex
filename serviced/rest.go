package serviced

import (
	"context"
	"net"

	"github.com/lamhai1401/rpc-ex/security"
	"google.golang.org/grpc"
)

// khai báo struct xây dựng RestService
type RestServiceImpl struct{}

// hàm Get RPC được xây dựng như sau
func (r *RestServiceImpl) Get(ctx context.Context, message *security.StringMessage) (*security.StringMessage, error) {
	return &security.StringMessage{Value: "Get hi:" + message.Value + "#"}, nil
}

// tương tự với hàm Post RPC được xây dựng với
func (r *RestServiceImpl) Post(ctx context.Context, message *security.StringMessage) (*security.StringMessage, error) {
	return &security.StringMessage{Value: "Post hi:" + message.Value + "@"}, nil
}

// hàm main của gRPC service
func MakeGrpcServer() *grpc.ClientConn {
	// khởi tạo một grpc Server mới
	grpcServer := grpc.NewServer()
	// grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
	// 	MaxConnectionIdle: 5 * time.Minute, // <--- This fixes it!
	// }))
	// register grpc Server với đối tượng xây dựng các hàm RPC
	security.RegisterRestServiceServer(grpcServer, new(RestServiceImpl))
	// listen gRPC Service trên port 5000, bỏ qua lỗi trả về nếu có
	lis, _ := net.Listen("tcp", ":5000")

	grpcServer.Serve(lis)

	// conn, err := grpc.DialContext(
	// 	context.Background(),
	// 	":5000",
	// 	grpc.WithInsecure(),
	// )
	// if err != nil {
	// 	log.Panicf("Failed to dial server %v", zap.Error(err))
	// }

	return nil
}
