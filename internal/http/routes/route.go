package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/jeevic/lego-demo/internal/http/controllers"
)

func InitRoute(e *gin.Engine) {

	//全局限速 针对各个url限速不在生效
	//ratelimiter.RateLimiter.AddWholeRateLimit(100)

	e.GET("/", controllers.HelloWorld)

	e.GET("/hello", controllers.Hello)

	e.GET("/grpcHello", controllers.GrpcHello)

	e.GET("/grpcHello1", controllers.GrpcHello1)

	//限流1000 qps
	//ratelimiter.RateLimiter.AddRateLimit("/", 1000)
	e.GET("/grpc", controllers.GrpcClient)
	//限流5000 qps
	//ratelimiter.RateLimiter.AddRateLimit("/grpc", 5000)

	e.GET("/Logg", controllers.Logg)

	e.GET("/Gos", controllers.Gos)
	e.GET("/getFeature", controllers.GetFeature)
	e.GET("/getFeatures", controllers.GetFeatures)
	e.GET("/runRecordRoute", controllers.RunRecordRoute)
	e.GET("/runRouteChat", controllers.RunRouteChat)
	e.GET("/Test", controllers.Test)
	e.GET("/Test2", controllers.Test2)

}
