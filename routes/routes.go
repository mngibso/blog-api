package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/handlers"
)

func SetupRoutes(router *gin.Engine) {
	// authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts
	// authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts
	// "foo":    "bar",
	// "austin": "1234",
	// "lena":   "hello2",
	// "manu":   "4321",
	// }))
	auth := router.Group("/v1", basicAuth())
	{
		auth.PUT("/post/:id", handlers.UpdatePost)
		auth.DELETE("/post/:id", handlers.DeletePost)
		auth.POST("/post", handlers.CreatePost)
		auth.PUT("/user/:username", handlers.UpdateUser)
		auth.DELETE("/user/:username", handlers.DeleteUser)
	}

	v1 := router.Group("/v1/")
	{
		v1.GET("/post", handlers.GetPost)
		v1.GET("/post/:id", handlers.GetPostById)
		v1.POST("/user", handlers.CreateUser)
		v1.GET("/user", handlers.GetUsers)
		v1.GET("/user/:username", handlers.GetUserByUserName)
	}
}
