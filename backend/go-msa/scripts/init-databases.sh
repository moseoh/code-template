#!/bin/sh

set -e

# PG* 환경 변수를 psql이 자동으로 사용하므로 인자를 넘길 필요가 없습니다.
psql -v ON_ERROR_STOP=1 <<-EOSQL
    -- auth_service 데이터베이스가 없으면 생성
    SELECT 'CREATE DATABASE auth_service'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth_service')\gexec

    -- business_service 데이터베이스가 없으면 생성
    SELECT 'CREATE DATABASE business_service'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'business_service')\gexec
EOSQL

echo "Databases are ready."