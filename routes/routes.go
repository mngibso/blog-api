package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/db"
	"github.com/mngibso/blog-api/handlers"
)

// SetupRoutes defines rest api routes and middleware on those routes
func SetupRoutes(router *gin.Engine, u db.UserStorer, p db.PostStorer) {

	userAPI := handlers.NewUserAPI(u)
	postAPI := handlers.NewPostAPI(p)
	// auth routes are protected by basic authentication using basicAuth middleware
	auth := router.Group("/v1", basicAuth(u))
	{
		auth.PUT("/post/:id", postAPI.UpdatePost)
		auth.DELETE("/post/:id", postAPI.DeletePost)
		auth.POST("/post", postAPI.CreatePost)
		auth.PUT("/user/:username", userAPI.UpdateUser)
		auth.DELETE("/user/:username", userAPI.DeleteUser)
	}

	v1 := router.Group("/v1")
	{
		v1.GET("/post", postAPI.GetPost)
		v1.GET("/post/:id", postAPI.GetPostById)
		v1.POST("/user", userAPI.CreateUser)
		v1.GET("/user", userAPI.GetUsers)
		v1.GET("/user/:username", userAPI.GetUserByUserName)
	}
}
