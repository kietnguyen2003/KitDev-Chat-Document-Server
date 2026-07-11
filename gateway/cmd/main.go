package main

import (
	gateway "gateway/internal"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config := gateway.LoadConfig()
	r := gin.New()
	rateLimiter := gateway.NewRateLimiter(
		config.RateLimitRequests,
		config.RateLimitWindow,
		config.RateLimitCleanupWindow,
	)

	gateway.CorsGuard(r)
	r.Use(rateLimiter.Middleware())

	authProxy := gateway.NewLoadBalancedProxy(config.ServerURLs).Handler()
	ragProxy := gateway.NewLoadBalancedProxy(config.RagURLs).Handler()

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authProxy)
		auth.POST("/sign-in", authProxy)
	}

	protected := r.Group("/api")
	protected.Use(gateway.AuthGuard())
	{
		protected.POST("/categories", authProxy)
		protected.GET("/categories", authProxy)

		protected.POST("/documents", authProxy)
		protected.DELETE("/documents", authProxy)
		protected.GET("/documents/:cateName", authProxy)
		protected.PUT("/documents/:id", authProxy)

		protected.POST("/ask/:language", ragProxy)
	}

	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal("Server down")
	}
}
