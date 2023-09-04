package main

import (
	"app-auth/pkg/database"
	"app-auth/web/handlers"
	"app-auth/web/middlewares"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitMongoClient()
	r := gin.Default()
	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)
	r.Use(middlewares.AuthMiddleware())
	r.POST("/admin/adduser", handlers.AddUser)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server could not start: %v", err)
	}
}
