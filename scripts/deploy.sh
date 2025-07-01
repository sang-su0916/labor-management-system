#!/bin/bash

# 노무관리 시스템 배포 스크립트

set -e

echo "🚀 노무관리 시스템 배포 시작..."

# 환경 변수 확인
if [ ! -f ".env" ]; then
    echo "❌ .env 파일이 없습니다. .env.example을 참조하여 생성해주세요."
    exit 1
fi

# Docker 및 Docker Compose 설치 확인
if ! command -v docker &> /dev/null; then
    echo "❌ Docker가 설치되지 않았습니다."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose가 설치되지 않았습니다."
    exit 1
fi

# 필요한 디렉토리 생성
echo "📁 필요한 디렉토리 생성..."
mkdir -p ./documents ./uploads ./logs ./data
mkdir -p ./nginx/ssl

# SSL 인증서 확인 (Let's Encrypt 사용 권장)
if [ ! -f "./nginx/ssl/fullchain.pem" ] || [ ! -f "./nginx/ssl/privkey.pem" ]; then
    echo "⚠️  SSL 인증서가 없습니다. 자체 서명 인증서를 생성합니다."
    echo "실제 운영환경에서는 Let's Encrypt 등을 사용하세요."
    
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout ./nginx/ssl/privkey.pem \
        -out ./nginx/ssl/fullchain.pem \
        -subj "/C=KR/ST=Seoul/L=Seoul/O=Company/CN=localhost"
fi

# 기존 컨테이너 중지 및 제거
echo "🛑 기존 컨테이너 중지..."
docker-compose down --remove-orphans

# 이미지 빌드
echo "🔨 Docker 이미지 빌드..."
docker-compose build --no-cache

# 데이터베이스 마이그레이션 (PostgreSQL 사용시)
echo "🗄️  데이터베이스 준비..."
docker-compose up -d postgres redis
sleep 10

# 애플리케이션 시작
echo "🚀 애플리케이션 시작..."
docker-compose up -d

# 헬스체크
echo "🔍 헬스체크 수행..."
sleep 30

if curl -f http://localhost:8080/ > /dev/null 2>&1; then
    echo "✅ 배포 성공! 애플리케이션이 정상적으로 실행 중입니다."
    echo "🌐 접속 URL: http://localhost (HTTP)"
    echo "🔒 접속 URL: https://localhost (HTTPS)"
else
    echo "❌ 헬스체크 실패. 로그를 확인해주세요."
    docker-compose logs labor-management
    exit 1
fi

# 상태 출력
echo "📊 컨테이너 상태:"
docker-compose ps

echo "✨ 배포 완료!"
echo ""
echo "📋 관리 명령어:"
echo "  - 로그 확인: docker-compose logs -f labor-management"
echo "  - 컨테이너 재시작: docker-compose restart"
echo "  - 중지: docker-compose down"
echo "  - DB 백업: docker-compose exec postgres pg_dump -U labor_user labor_management > backup.sql"