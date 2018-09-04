package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	router_v1 "github.com/chrootlogin/go-docstore/internal/api/v1/router"

	// auto load env files
	_ "github.com/joho/godotenv/autoload"
	"github.com/chrootlogin/go-docstore/internal/auth"
)

var port = ""

func main() {
	if len(port) == 0 {
		port = "8000"
	}

	initRouter()
}

func initRouter() {
	router := gin.Default()

	// Allow cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization", "*")
	corsConfig.AddAllowMethods("HEAD", "GET", "POST", "PUT", "DELETE")
	router.Use(cors.New(corsConfig))

	// Ping Route
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// user context
	am := auth.GetAuthMiddleware()
	{
		router.POST("/user/login", am.LoginHandler)

		api := router.Group("/api/")
		api.Use(am.MiddlewareFunc())
		{
			router_v1.InitRouter(api.Group("/v1/"))
		}
	}

	router.Run(":" + port)
}