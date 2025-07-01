# 멀티 스테이지 빌드를 사용한 최적화된 Dockerfile

# Build stage
FROM golang:1.24-alpine AS builder

# 필요한 패키지 설치
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 파일 복사 및 의존성 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# CGO 활성화하여 빌드 (SQLite 지원을 위해)
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o main ./cmd/server

# Production stage
FROM alpine:latest

# 필요한 패키지 설치
RUN apk --no-cache add ca-certificates sqlite tzdata

# 시간대 설정
ENV TZ=Asia/Seoul

# 사용자 생성
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 작업 디렉토리 생성
WORKDIR /app

# 빌드된 바이너리 복사
COPY --from=builder /app/main .

# 정적 파일 및 템플릿 복사
COPY --from=builder /app/web ./web
COPY --from=builder /app/database ./database

# 필요한 디렉토리 생성
RUN mkdir -p ./documents ./uploads ./logs

# 권한 설정
RUN chown -R appuser:appgroup /app

# 사용자 변경
USER appuser

# 포트 노출
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# 실행
CMD ["./main"]