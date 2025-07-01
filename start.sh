#!/bin/bash

# Railway 배포용 시작 스크립트

# 환경변수 설정
export GIN_MODE=${GIN_MODE:-release}
export PORT=${PORT:-8080}

# 필요한 디렉토리 생성
mkdir -p ./documents ./uploads ./logs

# 애플리케이션 실행
echo "🚀 노무관리 시스템 시작 (Railway 배포)"
echo "포트: $PORT"
echo "모드: $GIN_MODE"

./main