#!/bin/bash

# Railway 전용 시작 스크립트

echo "🚀 Railway에서 노무관리 시스템 시작"
echo "포트: $PORT"
echo "DATABASE_URL: ${DATABASE_URL:0:30}..."

# 필요한 디렉토리 생성
mkdir -p ./documents ./uploads ./logs

# 빌드 및 실행
CGO_ENABLED=1 go build -o main cmd/server/main.go
./main