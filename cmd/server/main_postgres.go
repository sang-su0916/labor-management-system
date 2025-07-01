// +build postgres

package main

import (
	"labor-management-system/database"
	"labor-management-system/internal/handlers"
	"labor-management-system/internal/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Set Gin mode
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	// Initialize database (PostgreSQL only)
	if err := database.InitDatabase(""); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDatabase()

	// Initialize Gin router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORSMiddleware())

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
		// Test database connection
		if err := database.DB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"service": "labor-management-system",
				"error": "database connection failed",
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "labor-management-system",
			"database": "connected",
			"timestamp": time.Now().UTC(),
		})
	})

	// Public routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "노무관리 시스템",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Authentication
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Employee management
			employees := protected.Group("/employees")
			{
				employees.GET("", handlers.GetEmployees)
				employees.POST("", middleware.RequireRole("admin", "hr"), handlers.CreateEmployee)
				employees.GET("/:id", handlers.GetEmployee)
				employees.PUT("/:id", middleware.RequireRole("admin", "hr"), handlers.UpdateEmployee)
				employees.DELETE("/:id", middleware.RequireRole("admin"), handlers.DeleteEmployee)
			}
		}
	}

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	// Bind to all interfaces for cloud deployment
	host := "0.0.0.0"
	addr := host + ":" + port

	log.Printf("Server starting on %s", addr)
	log.Fatal(r.Run(addr))
}