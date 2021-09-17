package main

import (
	"github.com/gin-gonic/gin"
	"learn/http/pb"
	"net/http"
)

func main() {
	ginTest := gin.Default()
	ginTest.GET("/protobuf", func(c *gin.Context) {
		data := &pb.User{
			Name: "张三",
			Age:  20,
		}
		c.ProtoBuf(http.StatusOK, data)
	})

	ginTest.Run(":8080")
}
