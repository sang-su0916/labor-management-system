# 🏢 노무관리 시스템 (Labor Management System)

중소기업을 위한 종합 노무관리 시스템입니다. 직원 관리, 급여 계산, 근태 관리, 휴가 관리 등의 기능을 제공합니다.

## ✨ 주요 기능

### 👥 직원 관리
- 직원 정보 등록, 수정, 삭제
- 부서 및 직급 관리
- 직원 검색 및 필터링

### 💰 급여 관리
- 자동 급여 계산 (4대보험, 소득세 포함)
- 급여명세서 PDF 생성
- 급여 이력 관리

### ⏰ 근태 관리
- 출퇴근 기록
- 근무시간 자동 계산
- 근태 현황 조회

### 🏖️ 휴가 관리
- 휴가 신청 및 승인
- 연차 잔여일수 관리
- 휴가 유형별 관리

### 📄 문서 관리
- 급여명세서 자동 생성
- 재직증명서 발급
- 근로계약서 생성

## 🛠️ 기술 스택

### Backend
- **언어**: Go 1.24
- **웹 프레임워크**: Gin
- **데이터베이스**: SQLite/PostgreSQL
- **인증**: JWT
- **PDF 생성**: gofpdf

### Frontend
- **HTML5**, **CSS3**, **JavaScript**
- **Bootstrap 5**
- **반응형 디자인**

### DevOps
- **컨테이너**: Docker & Docker Compose
- **웹서버**: Nginx
- **모니터링**: Prometheus & Grafana
- **CI/CD**: GitHub Actions

## 🚀 빠른 시작

### 배포 옵션

#### 1. Render 배포 (추천 ⭐)
[Render](https://render.com)는 무료 티어를 제공하며 Go 애플리케이션 배포에 최적화되어 있습니다.

**배포 단계:**
1. **GitHub에 코드 푸시**
   ```bash
   git add .
   git commit -m "Render 배포 준비"
   git push origin main
   ```

2. **Render 계정 생성**
   - [render.com](https://render.com)에서 계정 생성
   - GitHub 계정으로 로그인 가능

3. **새 웹 서비스 생성**
   - Dashboard에서 "New +" 클릭
   - "Web Service" 선택
   - GitHub 저장소 연결 (labor-management-system)

4. **배포 설정**
   - **Name**: `labor-management-system`
   - **Environment**: `Go`
   - **Build Command**: 
     ```bash
     go mod download && mkdir -p bin && CGO_ENABLED=1 go build -o bin/main cmd/server/main.go && chmod +x start-render.sh
     ```
   - **Start Command**: `./start-render.sh`
   - **Plan**: Free

5. **환경 변수 설정**
   - `PORT`: `10000`
   - `GIN_MODE`: `release`

6. **배포 완료!** 🎉
   - 자동으로 빌드 및 배포됩니다
   - `https://your-app-name.onrender.com`으로 접속 가능

#### 2. Railway 배포
```yaml
# railway.toml 사용
```

#### 3. 로컬 개발

### 사전 요구사항
- Docker & Docker Compose
- Git

### 1. 프로젝트 클론
```bash
git clone https://github.com/your-username/labor-management-system.git
cd labor-management-system
```

### 2. 환경 설정
```bash
cp .env.example .env
# .env 파일을 편집하여 환경에 맞게 설정
```

### 3. 배포
```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh
```

### 4. 접속
- HTTP: http://localhost
- HTTPS: https://localhost
- 기본 계정: admin / admin123

## 📋 상세 설치 가이드

### 개발 환경 설정

#### 1. Go 개발 환경
```bash
# Go 1.24 설치 필요
go mod download
go run cmd/server/main.go
```

#### 2. 데이터베이스 설정
```bash
# SQLite (개발용)
sqlite3 labor_management.db < database/schema.sql

# PostgreSQL (운영용)
createdb labor_management
psql labor_management < database/postgres_schema.sql
```

### 운영 환경 배포

#### 1. Docker Compose 사용
```bash
# 프로덕션 모드로 실행
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

#### 2. 개별 컨테이너 실행
```bash
# 애플리케이션만 실행
docker build -t labor-management .
docker run -p 8080:8080 labor-management
```

## 🔧 설정

### 환경 변수
주요 환경 변수들을 `.env` 파일에서 설정할 수 있습니다:

```bash
# 서버 설정
PORT=8080
GIN_MODE=release

# 데이터베이스
DB_TYPE=postgres
DB_HOST=localhost
DB_USER=labor_user
DB_PASSWORD=secure_password

# JWT 보안
JWT_SECRET=your_super_secret_key

# 회사 정보
COMPANY_NAME=귀하의 회사명
COMPANY_ADDRESS=회사 주소
```

### 보안 설정
- JWT 토큰 기반 인증
- HTTPS/SSL 지원
- Rate Limiting
- CORS 설정

## 📊 모니터링

### Prometheus & Grafana
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin123)

### 주요 메트릭
- 응답 시간
- 에러율
- 사용자 세션
- 데이터베이스 성능

## 🔒 보안

### 인증 및 권한
- JWT 기반 인증
- 역할 기반 접근 제어 (RBAC)
- 세션 타임아웃

### 데이터 보호
- 비밀번호 bcrypt 해싱
- HTTPS 강제 사용
- SQL Injection 방지
- XSS 방지

## 🧪 테스트

```bash
# 단위 테스트 실행
go test ./...

# 커버리지 확인
go test -cover ./...

# 통합 테스트
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📖 API 문서

### 인증
```bash
POST /api/auth/login
POST /api/auth/register
```

### 직원 관리
```bash
GET /api/employees
POST /api/employees
GET /api/employees/:id
PUT /api/employees/:id
DELETE /api/employees/:id
```

### 급여 관리
```bash
GET /api/payroll
POST /api/payroll
GET /api/payroll/:id
PUT /api/payroll/:id
DELETE /api/payroll/:id
```

전체 API 문서는 [API.md](./docs/API.md)를 참조하세요.

## 🤝 기여하기

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다. [LICENSE](LICENSE) 파일을 참조하세요.

## 📞 지원

- 이슈 리포트: [GitHub Issues](https://github.com/your-username/labor-management-system/issues)
- 이메일: support@yourcompany.com
- 문서: [Wiki](https://github.com/your-username/labor-management-system/wiki)

## 🗺️ 로드맵

- [ ] 모바일 앱 개발
- [ ] 전자결재 시스템
- [ ] 인사평가 모듈
- [ ] 교육 관리 시스템
- [ ] API 외부 연동

---

**Made with ❤️ in Korea**
