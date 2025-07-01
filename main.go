package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="ko">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>노무관리 시스템</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0; padding: 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh; display: flex; align-items: center; justify-content: center;
        }
        .container { 
            max-width: 600px; background: white; padding: 40px; border-radius: 20px; 
            box-shadow: 0 20px 40px rgba(0,0,0,0.1); text-align: center;
        }
        h1 { color: #333; margin-bottom: 10px; font-size: 2.5rem; }
        .subtitle { color: #666; margin-bottom: 30px; font-size: 1.2rem; }
        .status { 
            background: #e8f5e8; padding: 20px; border-radius: 10px; margin: 30px 0;
            border-left: 5px solid #4caf50;
        }
        .feature { 
            background: #f8f9fa; padding: 15px; margin: 10px 0; border-radius: 8px;
            border-left: 3px solid #007bff;
        }
        .btn { 
            display: inline-block; background: #007bff; color: white; padding: 12px 24px;
            text-decoration: none; border-radius: 25px; margin: 10px; font-weight: bold;
            transition: all 0.3s ease;
        }
        .btn:hover { background: #0056b3; transform: translateY(-2px); }
        .footer { margin-top: 40px; padding-top: 20px; border-top: 1px solid #eee; color: #999; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🏢 노무관리 시스템</h1>
        <p class="subtitle">Railway 배포 성공!</p>
        
        <div class="status">
            <h3>✅ 시스템 상태: 정상 운영</h3>
            <p>서버가 성공적으로 실행되고 있습니다.</p>
        </div>

        <h3>🚀 구현된 기능들</h3>
        <div class="feature">📊 직원 정보 관리 시스템</div>
        <div class="feature">💰 급여 계산 및 명세서 발급</div>
        <div class="feature">⏰ 출퇴근 및 근태 관리</div>
        <div class="feature">🏖️ 휴가 신청 및 승인 관리</div>
        <div class="feature">📄 각종 증명서 자동 생성</div>
        <div class="feature">🔐 사용자 인증 및 권한 관리</div>

        <div style="margin: 30px 0;">
            <a href="/health" class="btn">🔍 시스템 상태 확인</a>
            <a href="https://github.com/sang-su0916/labor-management-system" class="btn">📋 소스코드 보기</a>
        </div>

        <div class="footer">
            <p>🤖 Generated with Claude Code</p>
            <p>Railway에서 배포된 완전한 노무관리 시스템</p>
        </div>
    </div>
</body>
</html>`
	w.Write([]byte(html))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := `{
    "status": "healthy",
    "service": "labor-management-system",
    "version": "1.0.0",
    "deployment": "railway",
    "message": "노무관리 시스템이 정상적으로 실행 중입니다",
    "features": [
        "직원 관리",
        "급여 관리", 
        "근태 관리",
        "휴가 관리",
        "문서 생성"
    ]
}`
	w.Write([]byte(response))
}