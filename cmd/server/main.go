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

	// Log current working directory for debugging
	if cwd, err := os.Getwd(); err == nil {
		log.Printf("Current working directory: %s", cwd)
	}

	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "labor_management.db"
	}
	
	if err := database.InitDatabase(dbPath); err != nil {
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
	
	// Load HTML templates - try multiple paths with debugging
	templatePaths := []string{"web/templates/*", "./web/templates/*", "/app/web/templates/*", "templates/*"}
	templateFound := false
	
	// First, let's check what files exist
	log.Println("Checking for template directories...")
	for _, dir := range []string{"web/templates", "./web/templates", "/app/web/templates", "templates"} {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			log.Printf("Found directory: %s", dir)
			files, _ := os.ReadDir(dir)
			for _, file := range files {
				log.Printf("  - %s", file.Name())
			}
		}
	}
	
	for _, pattern := range templatePaths {
		matches, err := filepath.Glob(pattern)
		log.Printf("Trying pattern %s: %d matches, error: %v", pattern, len(matches), err)
		if len(matches) > 0 {
			r.LoadHTMLGlob(pattern)
			log.Printf("Templates loaded from: %s (found %d files)", pattern, len(matches))
			templateFound = true
			break
		}
	}
	if !templateFound {
		log.Println("ERROR: No template files found in any of the expected locations")
		log.Println("This will cause the application to show fallback text instead of the actual UI")
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
		// Check if templates are loaded
		if !templateFound {
			// More detailed error response
			c.String(http.StatusOK, "Template Loading Error\n\nThe application is running but cannot find the HTML templates.\n\nExpected locations:\n- web/templates/\n- ./web/templates/\n- /app/web/templates/\n\nPlease check the deployment logs for more information.")
			return
		}
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
				employees.POST("/with-contract", middleware.RequireRole("admin", "hr"), handlers.CreateEmployeeWithContract)
				employees.GET("/:id", handlers.GetEmployee)
				employees.PUT("/:id", middleware.RequireRole("admin", "hr"), handlers.UpdateEmployee)
				employees.DELETE("/:id", middleware.RequireRole("admin"), handlers.DeleteEmployee)
			}

			// Employment contracts
			contracts := protected.Group("/contracts")
			{
				contracts.GET("", handlers.GetContracts)
				contracts.POST("", middleware.RequireRole("admin", "hr"), handlers.CreateContract)
				contracts.POST("/with-employee", middleware.RequireRole("admin", "hr"), handlers.CreateContractWithEmployee)
				contracts.GET("/:id", handlers.GetContract)
				contracts.PUT("/:id", middleware.RequireRole("admin", "hr"), handlers.UpdateContract)
				contracts.DELETE("/:id", middleware.RequireRole("admin"), handlers.DeleteContract)
			}

			// Payroll
			payroll := protected.Group("/payroll")
			{
				payroll.GET("", handlers.GetPayrollRecords)
				payroll.POST("", middleware.RequireRole("admin", "hr"), handlers.CreatePayrollRecord)
				payroll.GET("/:id", handlers.GetPayrollRecord)
				payroll.PUT("/:id", middleware.RequireRole("admin", "hr"), handlers.UpdatePayrollRecord)
				payroll.DELETE("/:id", middleware.RequireRole("admin"), handlers.DeletePayrollRecord)
			}

			// Attendance
			attendance := protected.Group("/attendance")
			{
				attendance.GET("", handlers.GetAttendanceLogs)
				attendance.POST("/clock-in", handlers.ClockIn)
				attendance.POST("/clock-out", handlers.ClockOut)
				attendance.GET("/employee/:id", handlers.GetEmployeeAttendance)
			}

			// Leave requests
			leaves := protected.Group("/leaves")
			{
				leaves.GET("", handlers.GetLeaveRequests)
				leaves.POST("", handlers.CreateLeaveRequest)
				leaves.GET("/:id", handlers.GetLeaveRequest)
				leaves.PUT("/:id/approve", middleware.RequireRole("admin", "hr"), handlers.ApproveLeaveRequest)
				leaves.PUT("/:id/reject", middleware.RequireRole("admin", "hr"), handlers.RejectLeaveRequest)
			}

			// Documents
			documents := protected.Group("/documents")
			{
				documents.GET("/templates", handlers.GetDocumentTemplates)
				documents.POST("/generate/:type", handlers.GenerateDocument)
				documents.GET("/employee/:id", handlers.GetEmployeeDocuments)
			}

			// System settings
			settings := protected.Group("/settings")
			settings.Use(middleware.RequireRole("admin"))
			{
				settings.GET("", handlers.GetSystemSettings)
				settings.PUT("", handlers.UpdateSystemSettings)
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