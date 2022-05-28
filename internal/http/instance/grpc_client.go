package instance

import (
	"github.com/jeevic/lego/components/grpc/grpcclient"
)

var GrpcClient *grpcclient.GrpcClient
var GrpcClientPool *grpcclient.Pool

func GrpcClientInstance() {
	addr := "127.0.0.1:8501"
	b, err := grpcclient.NewClient(addr, grpcclient.NewOptions())
	if err != nil {
		panic("grp client error")
	}
	GrpcClient = b
}

func GrpcClientPoolInstance() {
	addr := "127.0.0.1:8501"

	pool, err := grpcclient.NewPool(addr, grpcclient.WithPoolCap(10))
	if err != nil {
		panic(err)
	}
	GrpcClientPool = pool

}
