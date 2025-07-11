# Go MSA 서비스 관리 Makefile

.PHONY: help build build-api-gateway build-auth-service build-business-service
.PHONY: run run-api-gateway run-auth-service run-business-service
.PHONY: test test-api-gateway test-auth-service test-business-service
.PHONY: tidy tidy-pkg tidy-api-gateway tidy-auth-service tidy-business-service
.PHONY: clean dev-init fmt lint status ps
.PHONY: stop stop-auth stop-business stop-gateway restart dev

# 기본 명령어 목록 표시 (자동 추출)
help: ## 기본:사용 가능한 명령어 목록 표시
	@echo "사용 가능한 명령어:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "; last_cat=""} {split($$2, parts, ":"); cat=parts[1]; desc=parts[2]; if(cat != last_cat) {print "\n\033[33m" cat ":\033[0m"; last_cat=cat} printf "  \033[36m%-25s\033[0m %s\n", $$1, desc}'

# ====================[ Build ]====================
build: ## 빌드:모든 서비스 빌드
	@echo "🔨 Go Workspace 빌드 중..."
	go work sync
	go build -o .bin/api-gateway ./api-gateway/cmd/api
	go build -o .bin/auth-service ./services/auth/cmd/api
	go build -o .bin/business-service ./services/business/cmd/api
	@echo "✅ 모든 서비스 빌드 완료"

# API Gateway 빌드
build-api-gateway: ## 빌드:API Gateway 빌드
	@echo "🔨 API Gateway 빌드 중..."
	go build -o .bin/api-gateway ./api-gateway/cmd/api
	@echo "✅ API Gateway 빌드 완료"

# Auth Service 빌드
build-auth-service: ## 빌드:Auth Service 빌드
	@echo "🔨 Auth Service 빌드 중..."
	go build -o .bin/auth-service ./services/auth/cmd/api
	@echo "✅ Auth Service 빌드 완료"

# Business Service 빌드
build-business-service: ## 빌드:Business Service 빌드
	@echo "🔨 Business Service 빌드 중..."
	go build -o .bin/business-service ./services/business/cmd/api
	@echo "✅ Business Service 빌드 완료"

# 모든 서비스 실행 (백그라운드)
run: ## 실행:모든 서비스 실행 (백그라운드)
	@echo "🚀 모든 서비스 실행 중..."
	$(MAKE) run-auth-service &
	sleep 2
	$(MAKE) run-business-service &
	sleep 2
	$(MAKE) run-api-gateway &
	@echo "✅ 모든 서비스 실행 완료 (백그라운드)"

# API Gateway 실행
run-api-gateway: ## 실행:API Gateway 실행
	@echo "🚀 API Gateway 실행 중... (포트 8080)"
	cd api-gateway && export $$(cat .env | xargs) && go run cmd/api/main.go

# Auth Service 실행
run-auth-service: ## 실행:Auth Service 실행
	@echo "🚀 Auth Service 실행 중... (포트 8081)"
	cd services/auth && export $$(cat .env | xargs) && go run cmd/api/main.go

# Business Service 실행
run-business-service: ## 실행:Business Service 실행
	@echo "🚀 Business Service 실행 중... (포트 8082)"
	cd services/business && export $$(cat .env | xargs) && go run cmd/api/main.go

# 모든 서비스 종료
stop: ## 종료:모든 서비스 종료
	@echo "🛑 모든 서비스 종료 중..."
	-lsof -ti:8081 | xargs -r kill -9 || echo "Auth Service (8081) 없음"
	-lsof -ti:8082 | xargs -r kill -9 || echo "Business Service (8082) 없음"
	-lsof -ti:8080 | xargs -r kill -9 || echo "API Gateway (8080) 없음"
	@echo "✅ 모든 서비스 종료 완료"

# 모든 서비스 테스트
test: ## 테스트:모든 서비스 테스트
	@echo "🧪 Go Workspace 테스트 중..."
	go test ./...
	@echo "✅ 모든 서비스 테스트 완료"

# API Gateway 테스트
test-api-gateway: ## 테스트:API Gateway 테스트
	@echo "🧪 API Gateway 테스트 중..."
	go test ./api-gateway/...

# Auth Service 테스트
test-auth-service: ## 테스트:Auth Service 테스트
	@echo "🧪 Auth Service 테스트 중..."
	go test ./services/auth/...

# Business Service 테스트
test-business-service: ## 테스트:Business Service 테스트
	@echo "🧪 Business Service 테스트 중..."
	go test ./services/business/...

# 모든 서비스 의존성 정리
tidy: ## 정리:모든 서비스 의존성 정리
	@echo "📦 Go Workspace 의존성 정리 중..."
	go work sync
	go mod tidy -C ./pkg
	go mod tidy -C ./api-gateway
	go mod tidy -C ./services/auth
	go mod tidy -C ./services/business
	@echo "✅ 모든 서비스 의존성 정리 완료"

# PKG 의존성 정리
tidy-pkg: ## 정리:PKG 의존성 정리
	@echo "📦 pkg 의존성 정리 중..."
	go mod tidy -C ./pkg

# API Gateway 의존성 정리
tidy-api-gateway: ## 정리:API Gateway 의존성 정리
	@echo "📦 API Gateway 의존성 정리 중..."
	go mod tidy -C ./api-gateway

# Auth Service 의존성 정리
tidy-auth-service: ## 정리:Auth Service 의존성 정리
	@echo "📦 Auth Service 의존성 정리 중..."
	go mod tidy -C ./services/auth

# Business Service 의존성 정리
tidy-business-service: ## 정리:Business Service 의존성 정리
	@echo "📦 Business Service 의존성 정리 중..."
	go mod tidy -C ./services/business

# 빌드 아티팩트 정리
clean: ## 정리:빌드 아티팩트 정리
	@echo "🧹 빌드 아티팩트 정리 중..."
	rm -rf .bin/
	@echo "✅ 정리 완료"

# 개발 환경 초기화
dev-init: ## 개발:개발 환경 초기화
	@echo "🔧 개발 환경 초기화 중..."
	$(MAKE) tidy
	$(MAKE) build
	@echo "✅ 개발 환경 초기화 완료"

# 모든 서비스 형식 검사
fmt: ## 개발:모든 서비스 코드 포맷팅
	@echo "🎨 모든 서비스 코드 포맷팅 중..."
	cd api-gateway && go fmt ./...
	cd services/auth && go fmt ./...
	cd services/business && go fmt ./...
	@echo "✅ 모든 서비스 코드 포맷팅 완료"

# 모든 서비스 린트 검사
lint: ## 개발:모든 서비스 린트 검사
	@echo "🔍 모든 서비스 린트 검사 중..."
	cd api-gateway && go vet ./...
	cd services/auth && go vet ./...
	cd services/business && go vet ./...
	@echo "✅ 모든 서비스 린트 검사 완료"

# 서비스 상태 확인
status: ## 모니터링:서비스 상태 확인
	@echo "📊 서비스 상태 확인 중..."
	@echo "API Gateway (8080): $$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health || echo "DOWN")"
	@echo "Auth Service (8081): $$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8081/health || echo "DOWN")"
	@echo "Business Service (8082): $$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8082/health || echo "DOWN")"

# 실행 중인 서비스 프로세스 확인
ps: ## 모니터링:실행 중인 서비스 프로세스 확인
	@echo "🔍 실행 중인 서비스 포트 상태:"
	@echo "Auth Service (8081): $$(lsof -ti:8081 > /dev/null && echo "실행 중" || echo "없음")"
	@echo "Business Service (8082): $$(lsof -ti:8082 > /dev/null && echo "실행 중" || echo "없음")"
	@echo "API Gateway (8080): $$(lsof -ti:8080 > /dev/null && echo "실행 중" || echo "없음")"

# 특정 서비스 종료
stop-auth: ## 종료:Auth Service 종료
	@echo "🛑 Auth Service 종료 중..."
	-lsof -ti:8081 | xargs -r kill -9 || echo "Auth Service (8081) 없음"
	@echo "✅ Auth Service 종료 완료"

stop-business: ## 종료:Business Service 종료
	@echo "🛑 Business Service 종료 중..."
	-lsof -ti:8082 | xargs -r kill -9 || echo "Business Service (8082) 없음"
	@echo "✅ Business Service 종료 완료"

stop-gateway: ## 종료:API Gateway 종료
	@echo "🛑 API Gateway 종료 중..."
	-lsof -ti:8080 | xargs -r kill -9 || echo "API Gateway (8080) 없음"
	@echo "✅ API Gateway 종료 완료"

# 서비스 재시작
restart: ## 재시작:모든 서비스 재시작
	@echo "🔄 모든 서비스 재시작 중..."
	$(MAKE) stop
	sleep 2
	$(MAKE) run
	@echo "✅ 모든 서비스 재시작 완료"

# 개발 워크플로우 (포맷팅 -> 린트 -> 테스트 -> 빌드)
dev: ## 개발:개발 워크플로우 실행 (포맷팅->린트->테스트->빌드)
	@echo "🚀 개발 워크플로우 실행 중..."
	$(MAKE) fmt
	$(MAKE) lint
	$(MAKE) test
	$(MAKE) build
	@echo "✅ 개발 워크플로우 완료"

# Flyway
flyway-migrate-auth-service: ## flyway:Auth Service Flyway Migrate
	@echo "🔨 Auth Service migration 중..."
	cd services/auth && export $$(cat .env | xargs) && flyway -locations="filesystem:./db/migration" migrate
	@echo "✅ Auth Service migration 완료"

flyway-migrate-business-service: ## flyway:Business Service Flyway Migrate
	@echo "🔨 Business Service migration 중..."
	cd services/business && export $$(cat .env | xargs) && flyway -locations="filesystem:./db/migration" migrate
	@echo "✅ Business Service migration 완료"

flyway-info-auth-service: ## flyway:Auth Service Flyway Info
	@echo
	cd services/auth && export $$(cat .env | xargs) && flyway -locations="filesystem:./db/migration" info
	@echo

flyway-info-business-service: ## flyway:Business Service Flyway Info
	@echo
	cd services/business && export $$(cat .env | xargs) && flyway -locations="filesystem:./db/migration" info
	@echo