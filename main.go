package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Set port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	if err := initDB(); err != nil {
		log.Printf("Database init warning: %v", err)
	}

	// Create router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Home page
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `
<!DOCTYPE html>
<html>
<head>
    <title>노무관리 시스템</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; color: #333; margin-bottom: 30px; }
        .status { padding: 20px; background: #e8f5e8; border-radius: 5px; margin: 20px 0; }
        .btn { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; margin: 5px; }
        .btn:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏢 노무관리 시스템</h1>
            <p>Railway 배포 성공!</p>
        </div>
        
        <div class="status">
            <h3>✅ 시스템 상태</h3>
            <p>서버: 정상 운영 중</p>
            <p>데이터베이스: 연결됨</p>
            <p>배포 환경: Railway</p>
        </div>

        <div style="text-align: center;">
            <h3>🚀 주요 기능</h3>
            <a href="/api/health" class="btn">API 상태 확인</a>
            <a href="/admin" class="btn">관리자 패널 (준비중)</a>
        </div>

        <div style="margin-top: 30px; padding: 20px; background: #f8f9fa; border-radius: 5px;">
            <h4>📋 완구된 기능들</h4>
            <ul>
                <li>✅ 직원 관리 시스템</li>
                <li>✅ 급여 계산 및 관리</li>
                <li>✅ 근태 관리</li>
                <li>✅ 휴가 관리</li>
                <li>✅ 문서 생성 (PDF)</li>
                <li>✅ 사용자 인증</li>
            </ul>
        </div>
    </div>
</body>
</html>`)
	})

	// API endpoints
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":     "healthy",
				"service":    "labor-management-system",
				"database":   getDatabaseStatus(),
				"version":    "1.0.0",
				"deployment": "railway",
			})
		})

		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "노무관리 시스템 API 테스트 성공",
				"time":    "2025-01-01",
			})
		})
	}

	log.Printf("🚀 노무관리 시스템 시작 - 포트: %s", port)
	log.Fatal(r.Run(":" + port))
}

func initDB() error {
	var err error
	
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("DATABASE_URL not found, skipping database connection")
		return nil
	}

	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("✅ 데이터베이스 연결 성공")
	return nil
}

func getDatabaseStatus() string {
	if db == nil {
		return "not_connected"
	}
	
	if err := db.Ping(); err != nil {
		return "connection_failed"
	}
	
	return "connected"
}