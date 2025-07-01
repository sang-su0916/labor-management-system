package main

import (
	"labor-management-system/database"
	"labor-management-system/internal/handlers"
	"labor-management-system/internal/middleware"
	"log"
	"net/http"
	"os"

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

	// Serve static files
	r.Static("/static", "./web/static")
	
	// Load HTML templates
	r.LoadHTMLGlob("web/templates/*")

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "labor-management-system",
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

			// Employment contracts
			contracts := protected.Group("/contracts")
			{
				contracts.GET("", handlers.GetContracts)
				contracts.POST("", middleware.RequireRole("admin", "hr"), handlers.CreateContract)
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
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(r.Run(":" + port))
}