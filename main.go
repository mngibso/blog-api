package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mngibso/blog-api/routes"
)

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
