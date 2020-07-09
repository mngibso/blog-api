package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/handlers"
)

// SetupRoutes defines rest api routes and middleware on those routes
func SetupRoutes(router *gin.Engine) {

	// auth routes are protected by basic authentication using basicAuth middleware
	auth := router.Group("/v1", basicAuth())
	{
		auth.PUT("/post/:id", handlers.UpdatePost)
		auth.DELETE("/post/:id", handlers.DeletePost)
		auth.POST("/post", handlers.CreatePost)
		auth.PUT("/user/:username", handlers.UpdateUser)
		auth.DELETE("/user/:username", handlers.DeleteUser)
	}

	v1 := router.Group("/v1")
	{
		v1.GET("/post", handlers.GetPost)
		v1.GET("/post/:id", handlers.GetPostById)
		v1.POST("/user", handlers.CreateUser)
		v1.GET("/user", handlers.GetUsers)
		v1.GET("/user/:username", handlers.GetUserByUserName)
	}
}
