package main

import (
	"context"
	"io"

	"github.com/lamhai1401/rpc-ex/model"
)

// // định nghĩa service struct
// type HelloService struct{}

// // định nghĩa hàm service Hello, quy tắc:
// // 1. Hàm service phải public (viết hoa)
// // 2. Có hai tham số trong hàm
// // 3. Tham số thứ hai phải kiểu con trỏ
// // 4. Phải trả về kiểu error

// func (p *HelloService) Hello(request *model.String, reply *model.String) error {
// 	// các hàm như .GetValue() đã được tạo ra trong file hello.pb.go
// 	reply.Value = "Hello, " + request.GetValue()

// 	// trả về error = nil nếu thành công
// 	return nil
// }

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *model.String,
) (*model.String, error) {
	reply := &model.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func (p *HelloServiceImpl) Channel(stream model.HelloService_ChannelServer) error {
	for {
		// Server nhận dữ liệu được gửi từ client
		// trong vòng lặp.
		args, err := stream.Recv()
		if err != nil {
			// Nếu gặp `io.EOF`, client stream sẽ đóng.
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &model.String{Value: "hello:" + args.GetValue()}

		// Dữ liệu trả về được  gửi đến client
		// thông qua stream và việc gửi nhận
		// dữ liệu stream hai chiều là hoàn toàn độc lập
		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}
