#!/bin/bash

# Render 배포용 시작 스크립트
set -e

echo "🚀 Render에서 노무관리 시스템 시작"
echo "포트: $PORT"

# 필요한 디렉토리 생성
mkdir -p ./documents ./uploads ./logs

# 환경 변수 설정
export GIN_MODE=release
export PORT=${PORT:-10000}

# 데이터베이스 URL이 있으면 PostgreSQL 사용, 없으면 SQLite 사용
if [ -n "$DATABASE_URL" ]; then
    echo "🗄️ PostgreSQL 데이터베이스 사용"
else
    echo "🗄️ SQLite 데이터베이스 사용"
    # SQLite 데이터베이스 초기화
    if [ ! -f "labor_management.db" ]; then
        echo "📊 데이터베이스 스키마 초기화 중..."
        if [ -f "database/schema.sql" ]; then
            sqlite3 labor_management.db < database/schema.sql
            echo "✅ 데이터베이스 스키마 적용 완료"
        fi
    fi
fi

echo "🎯 서버 시작 중..."
exec ./bin/main 