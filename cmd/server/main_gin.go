package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "labor-management-system",
			"timestamp": time.Now().UTC(),
			"version": "gin-test",
		})
	})

	// Home page
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Labor Management System with Gin",
			"status": "running",
			"timestamp": time.Now().UTC(),
		})
	})

	// Test API endpoint
	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API endpoint working",
			"framework": "gin",
		})
	})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	// Bind to all interfaces for cloud deployment
	host := "0.0.0.0"
	addr := host + ":" + port

	log.Printf("Gin server starting on %s", addr)
	log.Fatal(r.Run(addr))
}