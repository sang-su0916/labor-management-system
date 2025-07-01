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

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "labor-management-system",
			"timestamp": time.Now().UTC(),
			"version": "gin-test",
		})
	})

	// Home page with login
	r.GET("/", func(c *gin.Context) {
		html := `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>노무관리 시스템</title><style>body{font-family:Arial;margin:0;padding:20px;background:#f5f5f5}.container{max-width:400px;margin:100px auto;background:white;padding:30px;border-radius:8px;box-shadow:0 2px 10px rgba(0,0,0,0.1)}.title{text-align:center;color:#333;margin-bottom:30px}.form-group{margin-bottom:20px}label{display:block;margin-bottom:5px;color:#555}input{width:100%;padding:10px;border:1px solid #ddd;border-radius:4px;box-sizing:border-box}button{width:100%;padding:12px;background:#007bff;color:white;border:none;border-radius:4px;cursor:pointer;font-size:16px}button:hover{background:#0056b3}.message{margin-top:15px;padding:10px;border-radius:4px;text-align:center}.success{background:#d4edda;color:#155724;border:1px solid #c3e6cb}.error{background:#f8d7da;color:#721c24;border:1px solid #f5c6cb}</style></head><body><div class="container"><h1 class="title">노무관리 시스템</h1><form id="loginForm"><div class="form-group"><label>사용자명</label><input type="text" id="username" required></div><div class="form-group"><label>비밀번호</label><input type="password" id="password" required></div><button type="submit">로그인</button></form><div id="message"></div><div style="margin-top:20px;text-align:center;color:#666;font-size:14px">테스트 계정: admin / admin</div></div><script>document.getElementById('loginForm').addEventListener('submit',async(e)=>{e.preventDefault();const username=document.getElementById('username').value;const password=document.getElementById('password').value;const messageDiv=document.getElementById('message');try{const response=await fetch('/api/auth/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({username,password})});const data=await response.json();if(data.success){messageDiv.innerHTML='<div class="message success">로그인 성공! 환영합니다, '+data.user.username+'님</div>';setTimeout(()=>{window.location.href='/dashboard'},1000)}else{messageDiv.innerHTML='<div class="message error">'+data.message+'</div>'}}catch(error){messageDiv.innerHTML='<div class="message error">서버 연결 오류: '+error.message+'</div>'}});</script></body></html>`
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

	// Dashboard page
	r.GET("/dashboard", func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html lang="ko">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>대시보드 - 노무관리 시스템</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f8f9fa; }
        .header { background-color: #007bff; color: white; padding: 1rem 2rem; display: flex; justify-content: space-between; align-items: center; }
        .main { padding: 2rem; }
        .card { background: white; padding: 1.5rem; margin-bottom: 1rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1rem; }
        .btn { display: inline-block; padding: 0.5rem 1rem; background-color: #007bff; color: white; text-decoration: none; border-radius: 4px; margin: 0.5rem 0.5rem 0.5rem 0; }
        .btn:hover { background-color: #0056b3; }
        .logout { background-color: #dc3545; }
        .logout:hover { background-color: #c82333; }
    </style>
</head>
<body>
    <div class="header">
        <h1>노무관리 시스템</h1>
        <div>
            <span>환영합니다, admin님</span>
            <a href="/" class="btn logout">로그아웃</a>
        </div>
    </div>
    <div class="main">
        <div class="grid">
            <div class="card">
                <h3>직원 관리</h3>
                <p>직원 정보 등록, 수정, 조회</p>
                <a href="/employees" class="btn">직원 목록</a>
                <a href="#" class="btn">신규 등록</a>
            </div>
            <div class="card">
                <h3>근태 관리</h3>
                <p>출퇴근 기록 및 근무시간 관리</p>
                <a href="#" class="btn">근태 현황</a>
                <a href="#" class="btn">근무시간 조회</a>
            </div>
            <div class="card">
                <h3>급여 관리</h3>
                <p>급여 계산 및 지급 관리</p>
                <a href="#" class="btn">급여 계산</a>
                <a href="#" class="btn">급여 대장</a>
            </div>
            <div class="card">
                <h3>시스템 설정</h3>
                <p>시스템 환경 설정 및 관리</p>
                <a href="/api/test" class="btn">API 테스트</a>
                <a href="/health" class="btn">시스템 상태</a>
            </div>
        </div>
    </div>
</body>
</html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, html)
	})

	// Employee management APIs
	r.GET("/api/employees", func(c *gin.Context) {
		employees := []gin.H{
			{"id": 1, "name": "김철수", "position": "개발자", "department": "IT부", "salary": 5000000},
			{"id": 2, "name": "이영희", "position": "디자이너", "department": "디자인부", "salary": 4500000},
			{"id": 3, "name": "박민수", "position": "매니저", "department": "영업부", "salary": 5500000},
		}
		c.JSON(http.StatusOK, gin.H{"employees": employees})
	})

	// Employee list page
	r.GET("/employees", func(c *gin.Context) {
		html := `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>직원 목록</title><style>body{font-family:Arial;margin:0;padding:0;background:#f8f9fa}.header{background:#007bff;color:white;padding:1rem 2rem;display:flex;justify-content:space-between;align-items:center}.main{padding:2rem}.table{width:100%;background:white;border-radius:8px;overflow:hidden;box-shadow:0 2px 4px rgba(0,0,0,0.1)}th,td{padding:1rem;text-align:left;border-bottom:1px solid #eee}th{background:#f8f9fa}.btn{display:inline-block;padding:0.5rem 1rem;background:#007bff;color:white;text-decoration:none;border-radius:4px;margin:0.5rem 0}</style></head><body><div class="header"><h1>직원 목록</h1><a href="/dashboard" class="btn">대시보드</a></div><div class="main"><table class="table"><thead><tr><th>ID</th><th>이름</th><th>직책</th><th>부서</th><th>급여</th></tr></thead><tbody id="employeeList"></tbody></table></div><script>fetch('/api/employees').then(r=>r.json()).then(data=>{const tbody=document.getElementById('employeeList');data.employees.forEach(emp=>{tbody.innerHTML+=`<tr><td>${emp.id}</td><td>${emp.name}</td><td>${emp.position}</td><td>${emp.department}</td><td>${emp.salary.toLocaleString()}원</td></tr>`})});</script></body></html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, html)
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