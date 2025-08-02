# Kubernetes 애플리케이션 배포

홈랩 K3s 클러스터에 애플리케이션 배포 및 관리

## 개요

Helm과 kubectl을 사용하여 Kubernetes 클러스터에 다양한 애플리케이션을 배포하고 관리합니다.

## 사전 준비

### 필수 소프트웨어
```bash
# kubectl 및 helm 설치 확인
kubectl version --client
helm version --short

# 클러스터 연결 확인
kubectl get nodes
```

### kubeconfig 설정
```bash
# 01-server-init에서 생성된 kubeconfig 사용
export KUBECONFIG=~/.kube/homelab-config
```

## 배포 가능한 애플리케이션

각 애플리케이션별 세부 설정은 해당 디렉토리의 배포 스크립트와 설정 파일에서 확인하세요.

### SigNoz 모니터링 스택

홈랩 K3s 클러스터 및 호스트 메트릭 모니터링

#### 기술 스택
- **Helm**: 패키지 관리
- **OpenTelemetry**: 메트릭 수집
- **ClickHouse**: 데이터 저장소
- **Traefik**: Ingress 컨트롤러

#### 주요 기능
- 호스트 시스템 메트릭 수집 및 시각화
- Kubernetes 클러스터 메트릭 모니터링
- 애플리케이션 성능 모니터링 (APM)
- Traefik Ingress를 통한 웹 접속

#### 실행 방법
```bash
cd signoz
uv run deploy.py
```

#### 접속 정보
- **URL**: http://monitoring.moseoh.com
- **포트 포워딩**: `kubectl port-forward -n signoz svc/signoz 28080:8080`
- **기본 로그인**: admin@signoz.io / admin

#### 배포 구성
- SigNoz 애플리케이션 (Core + UI)
- K8s 인프라 메트릭 수집 (DaemonSet)
- Traefik Ingress 설정
- 자동 설정 파일: `signoz-values.yaml`, `signoz-k8s-values.yaml`
