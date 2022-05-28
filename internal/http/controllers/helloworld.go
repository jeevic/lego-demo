package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeevic/lego/components/godis"
	"github.com/jeevic/lego/components/grpc/grpcclient"
	"github.com/jeevic/lego/components/httplib"
	"github.com/jeevic/lego/components/log"
	"github.com/jeevic/lego/pkg/app"

	"github.com/jeevic/lego-demo/internal/http/instance"
	"github.com/jeevic/lego-demo/pb"
)

func HelloWorld(c *gin.Context) {
	//app.App.GetLogger().Infof("hello world %s", time.Now().String())
	req := httplib.Get("http://127.0.0.1:8500/hello?value=hello%20world!")
	res, _ := req.String()
	c.String(200, res)
}

func Hello(c *gin.Context) {
	//app.App.GetLogger().Infof("hello world %s", time.Now().String())
	c.String(200, "hello world!")
}

func GrpcHello(c *gin.Context) {
	// app.App.GetLogger().Infof("hello world  grpc %s", time.Now().String())
	cli := pb.NewHelloClient(instance.GrpcClient.GetConn())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := cli.SayHello(ctx, &pb.Param{Value: "hello world!"})
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("grpc requset server error: %s", err.Error()))
		return
	}
	c.String(200, r.GetValue())
}

func GrpcHello1(c *gin.Context) {
	// app.App.GetLogger().Infof("hello world  grpc %s", time.Now().String())
	client, _ := instance.GrpcClientPool.GetClient()
	cli := pb.NewHelloClient(client.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := cli.SayHello(ctx, &pb.Param{Value: "hello world!"})
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("grpc requset server error: %s", err.Error()))
		return
	}
	c.String(200, r.GetValue())
}

func GrpcClient(c *gin.Context) {
	// app.App.GetLogger().Infof("hello world  grpc %s", time.Now().String())
	addr := "127.0.0.1:8501"
	conn, _ := grpcclient.NewClient(addr, grpcclient.NewOptions())
	defer conn.Close()
	cli := pb.NewHelloClient(conn.GetConn())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := cli.SayHello(ctx, &pb.Param{Value: "this is client"})
	if err != nil {
		c.String(200, fmt.Sprintf("grpc requset server error: %s", err.Error()))
		return
	}
	c.String(200, fmt.Sprintf("grpc requset server success: %s", r.GetValue()))

}

func Gos(c *gin.Context) {
	zkAddr := []string{"10.126.173.11:2181", "10.126.173.12:2181", "10.126.173.13:2181"}
	zkDir := "/jodis/codis-caifeng_test"

	//初始化
	pool, err := godis.Create().SetZookeeperClient(zkAddr, zkDir, 3000).SetDb(5).SetPoolSize(10).Build()
	if err != nil {
		return
	}
	cli, _ := pool.GetClient()

	cmd := cli.Set("test_1", 10, 100*time.Second)
	str, err := cmd.Result()
	fmt.Println(str)

	//关闭
	pool.Close()
}

func Logg(c *gin.Context) {

	//默认调用 instance app
	app.App.GetLogger().Info("this is log")

	//调用默认app1
	l := log.GetLogger("app1")
	if l == nil {
		return
	}
	l.Error("this is log error")

}

func Test(c *gin.Context) {
	addr := c.DefaultQuery("addr", "127.0.0.1:8501")
	clientNum, _ := strconv.Atoi(c.DefaultQuery("c", "10"))
	workerNum, _ := strconv.Atoi(c.DefaultQuery("g", "10"))
	totalCount, _ := strconv.Atoi(c.DefaultQuery("n", "200000"))
	str := fmt.Sprintf("server addr: %s, totalCount: %d, multi client: %d, worker num: %d\n",
		addr,
		totalCount,
		clientNum,
		workerNum,
	)
	startTime := time.Now()
	handleMultiClient(addr, totalCount, clientNum, workerNum)
	costTime := float64(time.Now().Sub(startTime).Nanoseconds()) / float64(1000) / float64(1000) / float64(1000)

	qps := float64(totalCount) / costTime
	str = str + fmt.Sprintf("multi client: %d, qps is %.0f\n", clientNum, qps)
	c.String(http.StatusOK, str)
}

func Test2(c *gin.Context) {
	addr := c.DefaultQuery("addr", "127.0.0.1:8501")
	clientNum, _ := strconv.Atoi(c.DefaultQuery("c", "10"))
	workerNum, _ := strconv.Atoi(c.DefaultQuery("g", "10"))
	totalCount, _ := strconv.Atoi(c.DefaultQuery("n", "200000"))
	str := fmt.Sprintf("server addr: %s, totalCount: %d, multi client: %d, worker num: %d\n",
		addr,
		totalCount,
		clientNum,
		workerNum,
	)
	startTime := time.Now()
	handleMultiClient2(addr, totalCount, clientNum, workerNum)
	costTime := float64(time.Now().Sub(startTime).Nanoseconds()) / float64(1000) / float64(1000) / float64(1000)

	qps := float64(totalCount) / costTime
	str = str + fmt.Sprintf("multi client: %d, qps is %.0f\n", clientNum, qps)
	c.String(http.StatusOK, str)
}

func handleMultiClient(addr string, totalCount, clientNum, workerNum int) {
	var wg sync.WaitGroup
	pool, _ := grpcclient.NewPool(addr, grpcclient.WithPoolCap(int64(clientNum)))
	defer pool.Close()
	for index := 0; index < workerNum; index++ {
		wg.Add(1)
		go func() {
			roundCount := totalCount / workerNum
			startTime := time.Now()
			for idx := 0; idx < roundCount; idx++ {
				ct := time.Now()
				c, _ := pool.GetClient()
				app.App.GetLogger().Infof("get grpc client latency: %s", time.Now().Sub(ct))
				cli := pb.NewHelloClient(c.Conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r, err := cli.SayHello(ctx, &pb.Param{Value: "hello world!"})
				if err != nil {
					continue
				}
				if r.GetValue() != "hello world!" {
					panic("grpc error")
				}
			}
			costTime := float64(time.Now().Sub(startTime).Nanoseconds()) / float64(1000) / float64(1000)
			l := costTime / float64(roundCount)
			app.App.GetLogger().Infof("lattency:%.5fms", l)
			wg.Done()
		}()
	}
	wg.Wait()
}

func handleMultiClient2(addr string, totalCount, clientNum, workerNum int) {
	var wg sync.WaitGroup
	pool, _ := grpcclient.NewPool(addr, grpcclient.WithPoolCap(int64(clientNum)))
	defer pool.Close()
	for index := 0; index < workerNum; index++ {
		wg.Add(1)
		go func() {
			roundCount := totalCount / workerNum
			startTime := time.Now()
			c, _ := pool.GetClient()
			for idx := 0; idx < roundCount; idx++ {
				cli := pb.NewHelloClient(c.Conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r, err := cli.SayHello(ctx, &pb.Param{Value: "hello world!"})
				if err != nil {
					continue
				}
				if r.GetValue() != "hello world!" {
					panic("grpc error")
				}
			}
			costTime := float64(time.Now().Sub(startTime).Nanoseconds()) / float64(1000) / float64(1000)
			l := costTime / float64(roundCount)
			app.App.GetLogger().Infof("lattency:%.5fms", l)
			wg.Done()
		}()
	}
	wg.Wait()
}
