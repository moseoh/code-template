# =================================================================
# 1. Builder Stage: Go 소스 코드를 컴파일하는 스테이지
# =================================================================
FROM golang:1.24-alpine AS builder

# MODULE: 빌드할 서비스의 디렉터리 이름 (e.g., services/auth)
# 이 인자는 docker build 명령어에서 --build-arg 로 전달받습니다.
ARG MODULE

# 작업 디렉터리 설정
WORKDIR /app

# Go 모듈 캐싱을 위해 go.mod, go.sum, go.work 파일을 먼저 복사합니다.
# 이렇게 하면 소스 코드가 변경되어도 종속성이 바뀌지 않았다면 캐시를 활용해 빌드 속도가 빨라집니다.
COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY api-gateway/go.mod ./api-gateway/
COPY services/auth/go.mod services/auth/go.sum ./services/auth/
COPY services/business/go.mod services/business/go.sum ./services/business/

# go.work 파일을 기반으로 종속성을 다운로드합니다.
RUN go work sync
RUN go mod download

# 나머지 소스 코드를 전부 복사합니다.
COPY pkg/ pkg/
COPY ${MODULE}/ ${MODULE}/

# 서비스 애플리케이션을 빌드합니다.
# -o: 출력 파일 경로 지정
# -ldflags="-w -s": 디버깅 정보를 제거하여 바이너리 크기를 줄입니다.
# CGO_ENABLED=0: C 바인딩 없이 정적으로 링크된 바이너리를 생성하여 어떤 환경에서도 잘 동작하게 합니다.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/main \
    ./${MODULE}/cmd/api

# =================================================================
# 2. Final Stage: 실제 실행될 최종 이미지를 만드는 스테이지
# =================================================================
FROM alpine:latest

# 빌더 스테이지에서 컴파일된 실행 파일만 복사합니다.
COPY --from=builder /app/main /app/main

# (선택) 타임존 설정이 필요하다면 추가
# RUN apk --no-cache add tzdata

# (선택) 인증서 관련 파일이 필요하다면 추가
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 애플리케이션이 실행될 포트를 노출합니다. (각 서비스에 맞게 변경 가능)
EXPOSE 8080

# 컨테이너가 시작될 때 실행할 명령어
ENTRYPOINT ["/app/main"]