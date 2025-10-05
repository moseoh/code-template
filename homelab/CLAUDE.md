# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

이 프로젝트는 홈랩 환경을 자동화하는 Infrastructure as Code (IaC) 프로젝트입니다. Ansible을 기반으로 Ubuntu 서버의 초기 설정부터 K3s 클러스터 설치, Kubernetes 리소스 배포까지 전체 홈랩 환경을 구성합니다.

## Project Structure

프로젝트는 **도구별 + 계층별 구조**로 구성되어 있습니다:

```
homelab/
├── Justfile                    # 전체 워크플로우 관리 (기능별)
├── ansible/                    # Ansible 관련 모든 것
│   ├── Justfile               # Ansible 세부 작업
│   ├── playbooks/
│   │   ├── server/            # Ubuntu Server 관련
│   │   │   ├── ubuntu-setup.yml
│   │   │   ├── wakeonlan-setup.yml
│   │   │   └── nfs-server-setup.yml
│   │   ├── kubernetes/        # K8s 클러스터 설정
│   │   │   └── k3s-setup.yml
│   │   └── common/            # 공통 (알림 등)
│   │       └── slack-alert.yml
│   ├── inventory/
│   │   ├── hosts.yml
│   │   └── group_vars/all/vault.yml
│   ├── files/
│   ├── ansible.cfg
│   ├── pyproject.toml
│   └── .vault_pass
└── kubernetes/                 # Kubernetes 리소스 관련 모든 것
    ├── Justfile               # K8s 세부 작업
    ├── infrastructure/        # 인프라 컴포넌트
    │   └── nfs-provisioner/
    │       ├── Justfile
    │       └── values.yml
    └── applications/          # 애플리케이션
        ├── signoz/
        │   ├── Justfile
        │   └── signoz-values.yaml
        └── kubeflow/
            ├── Justfile
            └── manifests-1.10.2/
```

## Development Commands

### 🚀 기본 워크플로우 (루트에서 실행)

```bash
# 전체 홈랩 환경 구축 (순차 실행)
just setup-all

# 단계별 설정
just setup-base         # Ubuntu 초기 설정 + WakeOnLAN
just setup-storage      # NFS Server + Provisioner
just setup-k3s          # K3s 클러스터
just setup-monitoring   # SigNoz 모니터링
just setup-kubeflow     # Kubeflow

# 상태 확인
just status-all         # 전체 시스템 상태
just ping               # Ansible 연결 테스트
just k8s-status         # Kubernetes 리소스 상태
just helm-status        # Helm releases 상태

# Vault 관리
just init-vault         # Vault 패스워드 생성 (최초 1회)
just vault-edit         # Vault 편집
just vault-view         # Vault 확인
```

### 📂 Ansible 세부 작업 (ansible/ 디렉토리)

```bash
cd ansible

# 연결 테스트
just ping

# Server 플레이북
just run-ubuntu-setup   # Ubuntu 서버 초기 설정
just run-wakeonlan      # WakeOnLAN 설정
just run-nfs-server     # NFS 서버 설정

# Kubernetes 플레이북
just run-k3s            # K3s 클러스터 설치

# Common 플레이북
just run-slack-alert    # Slack 알림 전송

# 확인 명령어
just check-ubuntu-setup
just check-nfs-server
just check-k3s

# Vault 관리
just vault-edit         # Vault 파일 안전 편집
just vault-view         # Vault 파일 내용 확인
just vault-create-pass  # Vault 패스워드 파일 생성
```

### ☸️ Kubernetes 세부 작업 (kubernetes/ 디렉토리)

```bash
cd kubernetes

# 전체 상태 확인
just status-all
just helm-list

# 빠른 설치
just install-nfs        # NFS Provisioner
just install-signoz     # SigNoz
just install-kubeflow   # Kubeflow

# 개별 모듈 작업
cd infrastructure/nfs-provisioner && just install
cd applications/signoz && just install
cd applications/kubeflow && just install
```

## Architecture

### Core Design Principles

1. **도구별 분리**: Ansible과 Kubernetes 완전 분리
2. **중복 제거**: Ansible 환경 한 곳에서만 관리
3. **계층적 구조**: playbooks 내부에서 server/kubernetes/common으로 분류
4. **워크플로우 통합**: 루트 Justfile에서 기능별 워크플로우 정의

### Key Components

1. **루트 Justfile**: 전체 워크플로우 관리 (setup-all, setup-storage 등)
2. **ansible/**: 모든 Ansible 관련 파일 (playbooks, inventory, vault)
3. **kubernetes/**: 모든 K8s 리소스 (infrastructure, applications)

### Workflow Example: NFS 스토리지 설정

```bash
# 루트에서 실행 (권장)
just setup-storage

# 내부 동작:
# 1. ansible/playbooks/server/nfs-server-setup.yml 실행
# 2. 10초 대기 (NFS 서버 시작)
# 3. kubernetes/infrastructure/nfs-provisioner 배포
```

### Variable Management

- **호스트 변수**: `ansible/inventory/hosts.yml`에 IP와 사용자명만 정의
- **민감 변수**: `ansible/inventory/group_vars/all/vault.yml`에 암호화 저장
- **플레이북 변수**: 각 playbook 내부에 독립적으로 정의

## Configuration

### Vault Setup (필수)

```bash
# 1. Vault 패스워드 파일 생성
just init-vault

# 2. Vault 파일 편집 (최초 1회)
just vault-edit
```

### Required Vault Variables

```yaml
# 타겟 서버 정보
target_ip: "192.168.0.101"
target_username: "moseoh"
target_password: "your_password"

# 인증 정보
ansible_password: "your_password"
ansible_become_password: "your_password"

# 로컬 공개키 경로
local_public_key_path: "~/.ssh/id_rsa.pub"

# Slack 알림 설정 (선택사항)
slack_webhook_url: "https://hooks.slack.com/services/..."
```

### Prerequisites

- Python 3.12+
- uv (Python 패키지 관리자)
- just (명령어 실행기)
- kubectl (Kubernetes CLI)
- helm (Kubernetes 패키지 관리자)
- SSH 공개키 생성 (`ssh-keygen -t rsa`)
- 타겟 서버에 SSH 접근 권한

## Development Guidelines

### 작업 위치

- **전체 워크플로우**: 루트 디렉토리에서 `just setup-*` 실행
- **Ansible 세부 작업**: `ansible/` 디렉토리에서 작업
- **Kubernetes 세부 작업**: `kubernetes/` 디렉토리에서 작업

### 새로운 플레이북 추가

1. 적절한 디렉토리에 playbook 추가:
   - Server 관련: `ansible/playbooks/server/`
   - Kubernetes 관련: `ansible/playbooks/kubernetes/`
   - 공통: `ansible/playbooks/common/`

2. `ansible/Justfile`에 실행 명령어 추가

3. 필요시 루트 `Justfile`에 워크플로우 추가

### 새로운 Kubernetes 리소스 추가

1. 적절한 디렉토리에 리소스 추가:
   - 인프라: `kubernetes/infrastructure/`
   - 애플리케이션: `kubernetes/applications/`

2. 하위 디렉토리에 Justfile 생성

3. `kubernetes/Justfile`에 모듈 추가

4. 필요시 루트 `Justfile`에 워크플로우 추가

### 보안

- 민감한 정보는 반드시 Vault로 암호화 (`just vault-edit`)
- `.vault_pass` 파일은 절대 커밋하지 않음 (.gitignore에 포함)
- 플레이북별 독립적인 변수 관리

      IMPORTANT: this context may or may not be relevant to your tasks. You should not respond to this context unless it is highly relevant to your task.
