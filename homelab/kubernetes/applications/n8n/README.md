# n8n - 워크플로우 자동화 플랫폼

n8n은 오픈소스 워크플로우 자동화 도구로, 다양한 서비스를 연결하고 자동화할 수 있습니다.

## 📋 구성 요소

- **PostgreSQL**: n8n 데이터베이스
- **n8n**: 워크플로우 자동화 애플리케이션

## 🚀 설치

### 사전 준비 (필수!)

**중요**: 설치 전에 반드시 `.env` 파일을 생성하고 비밀번호를 설정해야 합니다.

```bash
cd kubernetes/applications/n8n

# 1. .env 파일 생성
just init-env

# 2. .env 파일 편집 (실제 비밀번호 입력)
vi .env

# 또는 직접 복사
cp .env.example .env
vi .env
```

### 루트에서 실행 (권장)

```bash
# 홈랩 루트 디렉토리에서
just setup-n8n
```

### kubernetes 디렉토리에서 실행

```bash
cd kubernetes
just install-n8n
```

### n8n 디렉토리에서 직접 실행

```bash
cd kubernetes/applications/n8n
just install
```

## 📊 상태 확인

```bash
# 전체 상태 확인
just status

# PostgreSQL 상태만 확인
just status-postgres

# n8n 상태만 확인
just status-n8n
```

## 🌐 접속

### Ingress 접속 (권장)

```bash
# Ingress를 통한 외부 접속
# 브라우저에서 http://n8n.moseoh.com 접속

# Ingress 상태 확인
kubectl get ingress -n n8n
```

**주의**: `n8n.moseoh.com`이 DNS 또는 `/etc/hosts`에 설정되어 있어야 합니다.

### 로컬 접속 (포트 포워딩)

```bash
# n8n 웹 UI 접속
just port-forward

# 브라우저에서 http://localhost:5678 접속
```

### PostgreSQL 접속 (디버깅용)

```bash
# PostgreSQL 포트 포워딩
just port-forward-postgres

# psql로 접속
psql -h localhost -p 5432 -U n8n n8n
```

## 📝 로그 확인

```bash
# n8n 로그
just logs

# PostgreSQL 로그
just logs-postgres

# 모든 Pod 로그
just logs-all
```

## 🔧 관리

### 재시작

```bash
# n8n Pod 재시작
just restart

# PostgreSQL Pod 재시작
just restart-postgres
```

### 재설치

```bash
# 전체 재설치 (데이터 삭제됨!)
just reinstall
```

### 삭제

```bash
# n8n 전체 삭제
just uninstall
```

## 🛠️ 문제 해결

### Pod 상세 정보 확인

```bash
just describe
```

### 이벤트 확인

```bash
just events
```

### Shell 접속

```bash
# n8n Pod 접속
just shell

# PostgreSQL Pod 접속
just shell-postgres
```

### Secret 확인

```bash
# PostgreSQL 비밀번호 확인
just show-secrets

# ConfigMap 확인
just show-configmap
```

## 📁 디렉토리 구조

```
n8n/
├── Justfile                       # n8n 관리 명령어
├── README.md                      # 이 문서
├── .env.example                   # 환경 변수 예시 (Git 커밋됨)
├── .env                           # 실제 환경 변수 (gitignore, 로컬에만 존재)
└── manifests/                     # Kubernetes manifest 파일
    ├── namespace.yaml
    ├── postgres-secret.yaml.example   # Secret 예시 (Git 커밋됨)
    ├── postgres-configmap.yaml
    ├── postgres-pvc.yaml
    ├── postgres-deployment.yaml
    ├── postgres-service.yaml
    ├── n8n-pvc.yaml
    ├── n8n-deployment.yaml        # Secret을 환경변수로 참조
    ├── n8n-service.yaml
    └── n8n-ingress.yaml           # Traefik Ingress (n8n.moseoh.com)
```

## 🔐 보안

### PostgreSQL 비밀번호 설정 (환경 변수 방식)

이 프로젝트는 **환경 변수 기반 Secret 관리**를 사용합니다.
실제 비밀번호는 Git에 커밋되지 않고, `.env` 파일에 저장됩니다.

#### 초기 설정 (설치 전 필수!)

```bash
cd kubernetes/applications/n8n

# 1. .env 파일 생성
just init-env

# 2. .env 파일 편집
vi .env
```

#### .env 파일 형식

```bash
POSTGRES_USER=n8n
POSTGRES_PASSWORD=your_secure_password_here       # 변경 필수!
POSTGRES_DB=n8n
POSTGRES_NON_ROOT_USER=n8n
POSTGRES_NON_ROOT_PASSWORD=your_secure_password_here  # 변경 필수!
```

#### Secret 재생성 (비밀번호 변경 시)

```bash
# .env 파일 편집
vi .env

# Secret 재생성
just create-secret

# Pod 재시작 (새 Secret 적용)
just restart
just restart-postgres
```

### Git 안전성

- ✅ `.env` 파일은 `.gitignore`에 등록되어 Git에 커밋되지 않습니다
- ✅ `.env.example` 파일만 Git에 커밋됩니다 (예시용)
- ✅ `manifests/postgres-secret.yaml.example` 파일만 Git에 커밋됩니다 (참고용)

### 환경변수 연결

n8n Deployment는 Secret을 환경변수로 자동 참조합니다:

```yaml
# n8n-deployment.yaml
env:
  - name: DB_POSTGRESDB_USER
    valueFrom:
      secretKeyRef:
        name: postgres-secret
        key: POSTGRES_NON_ROOT_USER
  - name: DB_POSTGRESDB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: postgres-secret
        key: POSTGRES_NON_ROOT_PASSWORD
```

### 팀 협업 시

```bash
# 1. 팀원이 리포지토리 클론 후
cd kubernetes/applications/n8n

# 2. .env 파일 생성 및 비밀번호 입력
just init-env
vi .env

# 3. 설치 (Secret 자동 생성)
just install
```

## 📚 참고 자료

- [n8n 공식 문서](https://docs.n8n.io/)
- [n8n GitHub](https://github.com/n8n-io/n8n)
- [n8n Hosting GitHub](https://github.com/n8n-io/n8n-hosting/tree/main/kubernetes)

## ⚠️ 주의사항

- n8n은 전문 지식이 필요한 Self-hosting 솔루션입니다.
- 데이터 백업을 정기적으로 수행하세요.
- PVC를 삭제하면 모든 워크플로우 데이터가 삭제됩니다.
