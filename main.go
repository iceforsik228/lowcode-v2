package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lowcode-v2/controllers"
	"lowcode-v2/initializers"
	"lowcode-v2/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)

	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	if err := r.Run(); err != nil {
		panic("Failed to start server on port 8080")
	}
}
