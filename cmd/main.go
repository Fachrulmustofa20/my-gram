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
		userRouter.DELETE("/:userId", handler.DeleteUser)
	}
	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.POST("/", handler.CreatePhoto)
		photoRouter.GET("/", handler.GetAllPhotos)
		photoRouter.PUT("/:photoId", middlewares.PhotoAuthorization(), handler.UpdatePhoto)
		photoRouter.DELETE("/:photoId", middlewares.PhotoAuthorization(), handler.DeletePhoto)
	}

	commentRouter := r.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.POST("/", handler.CreateComment)
		commentRouter.PUT("/:commentId", middlewares.CommentAuthorization(), handler.UpdateComment)
	}

	socialMediaRouter := r.Group("/socialmedias")
	{
		socialMediaRouter.Use(middlewares.Authentication())
		socialMediaRouter.POST("/", handler.CreateSocialMedia)
		socialMediaRouter.GET("/", handler.GetSocialMedia)
		socialMediaRouter.PUT("/:socialMediaId", middlewares.SocialMediaAuthorization(), handler.UpdateSocialMedia)
		socialMediaRouter.DELETE("/:socialMediaId", middlewares.SocialMediaAuthorization(), handler.DeleteSocialMedia)
	}
	return r
}
