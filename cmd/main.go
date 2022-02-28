package main

import (
	"fmt"
	"mygram/handler"
	"mygram/infra"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.InitDB()
	r := StartApp()

	fmt.Println("Server run with port :8080")
	r.Run(":8080")
}

func StartApp() *gin.Engine {
	r := gin.Default()
	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", handler.UserRegistration)
	}
	return r
}
