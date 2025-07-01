package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Serve static files - try multiple paths
	staticPaths := []string{"./web/static", "web/static", "/app/web/static"}
	staticFound := false
	for _, path := range staticPaths {
		if _, err := os.Stat(path); err == nil {
			r.Static("/static", path)
			log.Printf("Static files served from: %s", path)
			staticFound = true
			break
		}
	}
	if !staticFound {
		log.Println("Warning: No static files directory found")
	}
	
	// Load HTML templates - try multiple paths
	templatePaths := []string{"web/templates/*", "./web/templates/*", "/app/web/templates/*"}
	templateFound := false
	for _, pattern := range templatePaths {
		if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
			r.LoadHTMLGlob(pattern)
			log.Printf("Templates loaded from: %s", pattern)
			templateFound = true
			break
		}
	}
	if !templateFound {
		log.Println("Warning: No template files found")
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "labor-management-system",
			"database":  "pending",
			"timestamp": time.Now().UTC(),
			"version":   "static-files",
		})
	})

	// Home page
	r.GET("/", func(c *gin.Context) {
		if templateFound {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "노무관리 시스템",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":   "Labor Management System with Static Files",
				"status":    "running",
				"templates": "not found",
				"timestamp": time.Now().UTC(),
			})
		}
	})

	// Test API endpoint
	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":      "API endpoint working",
			"framework":    "gin",
			"static_files": staticFound,
			"templates":    templateFound,
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

	log.Printf("Static files server starting on %s", addr)
	log.Fatal(r.Run(addr))
}