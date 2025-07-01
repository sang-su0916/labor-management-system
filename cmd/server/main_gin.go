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

	// Home page with login form
	r.GET("/", func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html lang="ko">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>노무관리 시스템</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 400px; margin: 100px auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .title { text-align: center; color: #333; margin-bottom: 30px; }
        .form-group { margin-bottom: 20px; }
        label { display: block; margin-bottom: 5px; color: #555; }
        input { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 12px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; }
        button:hover { background-color: #0056b3; }
        .message { margin-top: 15px; padding: 10px; border-radius: 4px; text-align: center; }
        .success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="title">노무관리 시스템</h1>
        <form id="loginForm">
            <div class="form-group">
                <label for="username">사용자명</label>
                <input type="text" id="username" name="username" required>
            </div>
            <div class="form-group">
                <label for="password">비밀번호</label>
                <input type="password" id="password" name="password" required>
            </div>
            <button type="submit">로그인</button>
        </form>
        <div id="message"></div>
        <div style="margin-top: 20px; text-align: center; color: #666; font-size: 14px;">
            테스트 계정: admin / admin
        </div>
    </div>

    <script>
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const messageDiv = document.getElementById('message');
            
            try {
                const response = await fetch('/api/auth/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ username, password })
                });
                
                const data = await response.json();
                
                if (data.success) {
                    messageDiv.innerHTML = '<div class="message success">로그인 성공! 환영합니다, ' + data.user.username + '님</div>';
                    setTimeout(() => {
                        messageDiv.innerHTML += '<div class="message success">대시보드 준비 중...</div>';
                    }, 1000);
                } else {
                    messageDiv.innerHTML = '<div class="message error">' + data.message + '</div>';
                }
            } catch (error) {
                messageDiv.innerHTML = '<div class="message error">서버 연결 오류: ' + error.message + '</div>';
            }
        });
    </script>
</body>
</html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, html)
	})

	// Test API endpoint
	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API endpoint working",
			"framework": "gin",
		})
	})

	// Login API endpoint
	r.POST("/api/auth/login", func(c *gin.Context) {
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request format",
			})
			return
		}
		
		// 임시 로그인 로직 (admin/admin)
		if loginData.Username == "admin" && loginData.Password == "admin" {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Login successful",
				"token":   "temp_token_12345",
				"user": gin.H{
					"id":       1,
					"username": "admin",
					"role":     "admin",
				},
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid credentials",
			})
		}
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