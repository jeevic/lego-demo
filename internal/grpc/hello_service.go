package grpc

import (
	"context"

	"github.com/jeevic/lego/components/grpc/grpcserver"
	grpc_ratelimiter "github.com/jeevic/lego/components/grpc/grpcserver/grpc-ratelimiter"
	"github.com/jeevic/lego/pkg/app"

	"github.com/jeevic/lego-demo/pb"
)

func init() {
	grpc_ratelimiter.RateLimiter.AddRateLimit("/pb.Hello/SayHello", 1000000)
}

type HelloService struct {
	pb.HelloServer
}

func (h *HelloService) SayHello(ctx context.Context, in *pb.Param) (*pb.Param, error) {
	//var a = 10
	//var d = 1
	//var b = d - 1
	//_ = a / b

	requestid := grpcserver.FromContextRequestId(ctx, app.App.GetRequestId())
	app.App.GetLogger().Infof("requestId:%s", requestid)
	return &pb.Param{Value: in.Value}, nil
}
