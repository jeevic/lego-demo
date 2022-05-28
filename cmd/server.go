package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/jeevic/lego/components/crontab"
	"github.com/jeevic/lego/components/mongo"
	"github.com/jeevic/lego/pkg/app"
	"github.com/jeevic/lego/pkg/bootstarp"
	"github.com/spf13/cobra"

	"github.com/jeevic/lego-demo/internal/cron/task"
	"github.com/jeevic/lego-demo/internal/grpc"
	"github.com/jeevic/lego-demo/internal/http/instance"
	"github.com/jeevic/lego-demo/internal/http/middleware"
	"github.com/jeevic/lego-demo/internal/http/routes"
	"github.com/jeevic/lego-demo/internal/models"
	"github.com/jeevic/lego-demo/pb"
	routeGuide "github.com/jeevic/lego-demo/pb/route_guide"
)

var (
	// Used for flags.
	serverFlag struct {
		cfgFile string
		env     string
	}

	ServerCmd = &cobra.Command{
		Use:   "lego-demo",
		Short: "lego-demo",
		Long:  `lego-demo`,

		// server子命令的实现，启动应用服务
		Run: func(cmd *cobra.Command, args []string) {
			//设置配置文件
			app.App.SetCfgFile(serverFlag.cfgFile)
			//设置环境
			_ = app.App.SetEnv(serverFlag.env)

			//注册启动项
			bootstarp.RegisterInit(func() {
				//注册启动mongo 也可init启动
				models.InitMongo()
				instance.GrpcClientInstance()
				instance.GrpcClientPoolInstance()
			})
			//启动
			err := bootstarp.Init()
			if err != nil {
				panic("bootstarp error")
			}

			//注册http路由
			_ = bootstarp.RegisterHttpRoutes(routes.InitRoute)

			//注册自定义custom
			hs, _ := app.App.GetHttpServer()
			hs.SetMiddleware(middleware.HelloWorld())

			//注册定时任务
			crontab.AddTaskFunc(task.Cron)
			//注册grpc server
			gs, _ := app.App.GetGrpcServer()
			pb.RegisterHelloServer(gs.Server, &grpc.HelloService{})
			routeGuide.RegisterRouteGuideServer(gs.Server, grpc.NewRouteGuideServer())

			//注册信号
			bootstarp.RegisterSignalFunc(func(sig os.Signal) {
				switch sig {
				case syscall.SIGINT:
				case syscall.SIGHUP:

				}
			})
			//注册shutdown清理函数
			bootstarp.RegisterShutdown(func() {
				//停止定时任务
				crontab.Stop()

				//清理注册的mongo
				mongo.Reset()
			})

			//启动run
			bootstarp.Run()

		},
	}
)

// Execute executes the root command.
func Execute() {
	if err := ServerCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	ServerCmd.PersistentFlags().StringVarP(&serverFlag.cfgFile, "config", "c", "", "config file")
	ServerCmd.PersistentFlags().StringVarP(&serverFlag.env, "env", "e", "develop", "env:develop,test,release,prod")

}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if len(serverFlag.cfgFile) <= 0 {
		er("must set config file")
	}

	if len(serverFlag.env) <= 0 {
		er("must set env")
	}
}
