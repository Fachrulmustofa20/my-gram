package main

import (
	"fmt"
	"mygram/handler"
	"mygram/infra"
	"mygram/middlewares"

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
		userRouter.POST("/login", handler.UserLogin)

		userRouter.Use(middlewares.Authentication())
		userRouter.PUT("/:userId", handler.UpdateUser)
	}
	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.POST("/", handler.CreatePhoto)
	}
	return r
}
