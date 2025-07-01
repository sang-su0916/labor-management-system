#!/bin/bash

# Railway 전용 시작 스크립트
set -e

echo "🚀 Railway에서 노무관리 시스템 시작"
echo "포트: $PORT"
echo "현재 디렉토리: $(pwd)"
echo "파일 목록:"
ls -la

# 필요한 디렉토리 생성
mkdir -p ./documents ./uploads ./logs ./web/static ./web/templates

# Go 모듈 다운로드
echo "📦 Go 모듈 다운로드 중..."
go mod download

# 데이터베이스 스키마 초기화 (SQLite용)
if [ ! -f "labor_management.db" ]; then
    echo "🗄️ 데이터베이스 초기화 중..."
    if [ -f "database/schema.sql" ]; then
        sqlite3 labor_management.db < database/schema.sql
        echo "✅ 데이터베이스 스키마 적용 완료"
    fi
fi

# 빌드 및 실행
echo "🔨 애플리케이션 빌드 중..."
CGO_ENABLED=1 go build -o main cmd/server/main.go

echo "🚀 애플리케이션 시작..."
./main