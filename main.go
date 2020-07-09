package main

import (
	"github.com/mngibso/blog-api/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/routes"
)

func main() {
	r := gin.Default()
	db.InitializeMongoDB()
	u := db.NewUserStore()
	p := db.NewPostStore()
	routes.SetupRoutes(r, u, p)
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
