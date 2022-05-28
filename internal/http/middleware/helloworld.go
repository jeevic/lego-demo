package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func HelloWorld() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("hello world!")

		c.Next()
	}
}
