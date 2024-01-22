package main

import (
	"jwtgin/controllers"
	"jwtgin/middlewares"
	"jwtgin/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	models.ConnectDatabase()
	apiRt := router.Group("api")
	apiRt.POST("/register", controllers.Register)
	apiRt.POST("/login", controllers.Login)
	protected := router.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	router.Run(":8080")
}
