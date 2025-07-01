#!/bin/bash

# Render 배포용 시작 스크립트
set -e

echo "🚀 Render에서 노무관리 시스템 시작"
echo "포트: $PORT"
echo "현재 작업 디렉토리: $(pwd)"
echo "파일 목록:"
ls -la

# 필요한 디렉토리 생성
mkdir -p ./documents ./uploads ./logs

# 환경 변수 설정
export GIN_MODE=release
export PORT=${PORT:-10000}

# 데이터베이스 URL이 있으면 PostgreSQL 사용, 없으면 SQLite 사용
if [ -n "$DATABASE_URL" ]; then
    echo "🗄️ PostgreSQL 데이터베이스 사용"
    echo "📊 DATABASE_URL: ${DATABASE_URL}"
    # PostgreSQL 연결 테스트 (선택사항)
    # psql "$DATABASE_URL" -c "SELECT 1;" || echo "⚠️ PostgreSQL 연결 확인 불가 - 애플리케이션에서 처리됩니다"
else
    echo "🗄️ SQLite 데이터베이스 사용"
    # SQLite 데이터베이스 초기화
    if [ ! -f "labor_management.db" ]; then
        echo "📊 데이터베이스 스키마 초기화 중..."
        if [ -f "database/schema.sql" ]; then
            sqlite3 labor_management.db < database/schema.sql
            echo "✅ 데이터베이스 스키마 적용 완료"
        else
            echo "⚠️ SQLite 스키마 파일을 찾을 수 없습니다 - 애플리케이션에서 처리됩니다"
        fi
    else
        echo "✅ 기존 SQLite 데이터베이스 사용"
    fi
fi

echo "🎯 서버 시작 중..."
exec ./main 