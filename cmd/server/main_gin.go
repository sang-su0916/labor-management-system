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
			c.Header("Content-Type", "application/json")
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
			c.Header("Content-Type", "application/json")
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
                <a href="/contract" class="btn">근로계약서</a>
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
                <a href="/payroll" class="btn">급여명세서</a>
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

	// 근로계약서 작성 페이지
	r.GET("/contract", func(c *gin.Context) {
		html := `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>근로계약서 작성</title><style>body{font-family:Arial;margin:0;padding:0;background:#f8f9fa}.header{background:#007bff;color:white;padding:1rem 2rem;display:flex;justify-content:space-between;align-items:center}.main{padding:2rem;max-width:800px;margin:0 auto}.form{background:white;padding:2rem;border-radius:8px;box-shadow:0 2px 4px rgba(0,0,0,0.1)}.form-group{margin-bottom:1rem}label{display:block;margin-bottom:0.5rem;font-weight:bold}input,select,textarea{width:100%;padding:0.5rem;border:1px solid #ddd;border-radius:4px;box-sizing:border-box}textarea{height:100px}.btn{display:inline-block;padding:0.75rem 1.5rem;background:#007bff;color:white;text-decoration:none;border:none;border-radius:4px;cursor:pointer;margin:0.5rem 0.5rem 0.5rem 0}.btn:hover{background:#0056b3}.btn-secondary{background:#6c757d}.btn-secondary:hover{background:#545b62}</style></head><body><div class="header"><h1>근로계약서 작성</h1><a href="/dashboard" class="btn btn-secondary">대시보드</a></div><div class="main"><div class="form"><h2>근로계약서</h2><form id="contractForm"><div class="form-group"><label>직원명</label><input type="text" id="employeeName" required></div><div class="form-group"><label>주민등록번호</label><input type="text" id="idNumber" required></div><div class="form-group"><label>근무부서</label><input type="text" id="department" required></div><div class="form-group"><label>직책</label><input type="text" id="position" required></div><div class="form-group"><label>근무시작일</label><input type="date" id="startDate" required></div><div class="form-group"><label>월 급여액</label><input type="number" id="salary" required></div><div class="form-group"><label>근무시간</label><select id="workHours"><option value="9-18">09:00 - 18:00</option><option value="8-17">08:00 - 17:00</option><option value="10-19">10:00 - 19:00</option></select></div><div class="form-group"><label>특이사항</label><textarea id="notes" placeholder="기타 근로조건 및 특약사항"></textarea></div><button type="submit" class="btn">계약서 생성</button><button type="button" class="btn btn-secondary" onclick="printContract()">인쇄</button></form></div></div><script>function printContract(){const data={name:document.getElementById('employeeName').value,idNumber:document.getElementById('idNumber').value,department:document.getElementById('department').value,position:document.getElementById('position').value,startDate:document.getElementById('startDate').value,salary:document.getElementById('salary').value,workHours:document.getElementById('workHours').value,notes:document.getElementById('notes').value};if(!data.name){alert('직원명을 입력해주세요');return;}const printWindow=window.open('','_blank');printWindow.document.write(\`<html><head><title>근로계약서</title><style>body{font-family:serif;padding:2rem;line-height:1.6}.header{text-align:center;margin-bottom:2rem}.content{margin:1rem 0}.signature{margin-top:3rem;display:flex;justify-content:space-between}</style></head><body><div class="header"><h1>근 로 계 약 서</h1></div><div class="content"><p><strong>직원명:</strong> \${data.name}</p><p><strong>주민등록번호:</strong> \${data.idNumber}</p><p><strong>근무부서:</strong> \${data.department}</p><p><strong>직책:</strong> \${data.position}</p><p><strong>근무시작일:</strong> \${data.startDate}</p><p><strong>월 급여액:</strong> \${parseInt(data.salary).toLocaleString()}원</p><p><strong>근무시간:</strong> \${data.workHours}</p><p><strong>특이사항:</strong> \${data.notes}</p></div><div class="signature"><div>직원 서명: _______________</div><div>회사 서명: _______________</div></div><div style="text-align:center;margin-top:2rem">작성일: \${new Date().toLocaleDateString()}</div></body></html>\`);printWindow.document.close();printWindow.print();}</script></body></html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, html)
	})

	// 급여명세서 페이지
	r.GET("/payroll", func(c *gin.Context) {
		html := `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>급여명세서</title><style>body{font-family:Arial;margin:0;padding:0;background:#f8f9fa}.header{background:#007bff;color:white;padding:1rem 2rem;display:flex;justify-content:space-between;align-items:center}.main{padding:2rem;max-width:800px;margin:0 auto}.form{background:white;padding:2rem;border-radius:8px;box-shadow:0 2px 4px rgba(0,0,0,0.1)}.form-group{margin-bottom:1rem}label{display:block;margin-bottom:0.5rem;font-weight:bold}input,select{width:100%;padding:0.5rem;border:1px solid #ddd;border-radius:4px;box-sizing:border-box}.btn{display:inline-block;padding:0.75rem 1.5rem;background:#007bff;color:white;text-decoration:none;border:none;border-radius:4px;cursor:pointer;margin:0.5rem 0.5rem 0.5rem 0}.btn:hover{background:#0056b3}.btn-secondary{background:#6c757d}.btn-secondary:hover{background:#545b62}.payroll-table{width:100%;margin-top:1rem;border-collapse:collapse}.payroll-table th,.payroll-table td{padding:0.5rem;border:1px solid #ddd;text-align:right}.payroll-table th{background:#f8f9fa;text-align:center}</style></head><body><div class="header"><h1>급여명세서</h1><a href="/dashboard" class="btn btn-secondary">대시보드</a></div><div class="main"><div class="form"><h2>급여 계산</h2><form id="payrollForm"><div class="form-group"><label>직원명</label><select id="employee"><option value="">직원 선택</option><option value="김철수,5000000">김철수 (500만원)</option><option value="이영희,4500000">이영희 (450만원)</option><option value="박민수,5500000">박민수 (550만원)</option></select></div><div class="form-group"><label>급여년월</label><input type="month" id="payMonth" required></div><div class="form-group"><label>기본급</label><input type="number" id="baseSalary" readonly></div><div class="form-group"><label>연장수당</label><input type="number" id="overtime" value="0"></div><div class="form-group"><label>특별수당</label><input type="number" id="bonus" value="0"></div><button type="button" class="btn" onclick="calculatePay()">급여 계산</button><button type="button" class="btn btn-secondary" onclick="printPayroll()">명세서 출력</button></form><div id="payrollResult" style="display:none"><h3>급여 명세</h3><table class="payroll-table"><tr><th>항목</th><th>금액</th></tr><tr><td>기본급</td><td id="resultBase">0원</td></tr><tr><td>연장수당</td><td id="resultOvertime">0원</td></tr><tr><td>특별수당</td><td id="resultBonus">0원</td></tr><tr><td>총 지급액</td><td id="resultGross">0원</td></tr><tr><td>소득세</td><td id="resultIncomeTax">0원</td></tr><tr><td>국민연금</td><td id="resultPension">0원</td></tr><tr><td>건강보험</td><td id="resultHealth">0원</td></tr><tr><td>고용보험</td><td id="resultEmployment">0원</td></tr><tr><td>총 공제액</td><td id="resultDeduction">0원</td></tr><tr style="font-weight:bold;background:#e3f2fd"><td>실수령액</td><td id="resultNet">0원</td></tr></table></div></div></div><script>document.getElementById('employee').addEventListener('change',function(){const value=this.value;if(value){const[name,salary]=value.split(',');document.getElementById('baseSalary').value=salary;}});function calculatePay(){const baseSalary=parseInt(document.getElementById('baseSalary').value)||0;const overtime=parseInt(document.getElementById('overtime').value)||0;const bonus=parseInt(document.getElementById('bonus').value)||0;const gross=baseSalary+overtime+bonus;const incomeTax=Math.floor(gross*0.033);const pension=Math.floor(gross*0.045);const health=Math.floor(gross*0.0335);const employment=Math.floor(gross*0.008);const totalDeduction=incomeTax+pension+health+employment;const net=gross-totalDeduction;document.getElementById('resultBase').textContent=baseSalary.toLocaleString()+'원';document.getElementById('resultOvertime').textContent=overtime.toLocaleString()+'원';document.getElementById('resultBonus').textContent=bonus.toLocaleString()+'원';document.getElementById('resultGross').textContent=gross.toLocaleString()+'원';document.getElementById('resultIncomeTax').textContent=incomeTax.toLocaleString()+'원';document.getElementById('resultPension').textContent=pension.toLocaleString()+'원';document.getElementById('resultHealth').textContent=health.toLocaleString()+'원';document.getElementById('resultEmployment').textContent=employment.toLocaleString()+'원';document.getElementById('resultDeduction').textContent=totalDeduction.toLocaleString()+'원';document.getElementById('resultNet').textContent=net.toLocaleString()+'원';document.getElementById('payrollResult').style.display='block';}function printPayroll(){const employee=document.getElementById('employee').value;if(!employee){alert('직원을 선택해주세요');return;}const[name]=employee.split(',');const month=document.getElementById('payMonth').value;const printWindow=window.open('','_blank');printWindow.document.write(\`<html><head><title>급여명세서</title><style>body{font-family:serif;padding:2rem;line-height:1.6}.header{text-align:center;margin-bottom:2rem}.payroll-table{width:100%;border-collapse:collapse;margin:1rem 0}.payroll-table th,.payroll-table td{padding:0.5rem;border:1px solid #000;text-align:center}.payroll-table th{background:#f0f0f0}</style></head><body><div class="header"><h1>급 여 명 세 서</h1><p>지급년월: \${month}</p><p>성명: \${name}</p></div><table class="payroll-table"><tr><th>항목</th><th>금액</th></tr><tr><td>기본급</td><td>\${document.getElementById('resultBase').textContent}</td></tr><tr><td>연장수당</td><td>\${document.getElementById('resultOvertime').textContent}</td></tr><tr><td>특별수당</td><td>\${document.getElementById('resultBonus').textContent}</td></tr><tr><td>총 지급액</td><td>\${document.getElementById('resultGross').textContent}</td></tr><tr><td>소득세</td><td>\${document.getElementById('resultIncomeTax').textContent}</td></tr><tr><td>국민연금</td><td>\${document.getElementById('resultPension').textContent}</td></tr><tr><td>건강보험</td><td>\${document.getElementById('resultHealth').textContent}</td></tr><tr><td>고용보험</td><td>\${document.getElementById('resultEmployment').textContent}</td></tr><tr><td>총 공제액</td><td>\${document.getElementById('resultDeduction').textContent}</td></tr><tr style="font-weight:bold"><td>실수령액</td><td>\${document.getElementById('resultNet').textContent}</td></tr></table><div style="text-align:center;margin-top:2rem">발행일: \${new Date().toLocaleDateString()}</div></body></html>\`);printWindow.document.close();printWindow.print();}</script></body></html>`
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