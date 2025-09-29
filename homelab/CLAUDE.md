# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

이 프로젝트는 홈랩 환경을 자동화하는 Infrastructure as Code (IaC) 프로젝트입니다. Ansible을 기반으로 Ubuntu 서버의 초기 설정부터 K3s 클러스터 설치까지 전체 홈랩 환경을 단계별로 구성합니다.

## Project Structure

프로젝트는 단계별 구조로 구성되어 있습니다:

```
01-server-init/     # Ansible 기반 서버 초기 설정 및 K3s 설치
02-kubernetes/      # Kubernetes 리소스 정의 (향후 확장)
```

## Development Commands

### 환경 설정
```bash
# Python 가상환경 생성 및 활성화 (uv 사용)
cd 01-server-init
uv venv
source .venv/bin/activate

# 의존성 설치
uv pip install -r pyproject.toml

# Vault 패스워드 설정
echo "your_vault_password" > .vault_pass
chmod 600 .vault_pass
```

### Just 명령어 (권장)
```bash
cd 01-server-init

# 연결 테스트
just ping

# 실행 명령어 (run- prefix)
just run-ubuntu-setup    # Ubuntu 서버 초기 설정
just run-wakeonlan       # WakeOnLAN 설정
just run-nfs-server      # NFS 서버 설정
just run-k3s             # K3s 클러스터 설치
just run-all             # 전체 설치 (순차 실행)

# 확인 명령어 (check- prefix)
just check-ubuntu-setup  # Ubuntu 설정 변경사항 확인
just check-all           # 모든 플레이북 확인

# Vault 관리
just vault-edit          # Vault 파일 안전 편집
just vault-view          # Vault 파일 내용 확인
just vault-create-pass   # Vault 패스워드 파일 생성

# 상태 확인
just inventory-list      # 인벤토리 목록
```

## Architecture

### Core Components

1. **Ansible 표준 구조**: `inventory/group_vars/all/vault.yml`을 통한 Vault 기반 변수 관리
2. **Just 기반 실행**: Justfile을 통한 명령어 표준화 및 그룹화
3. **uv 통합**: 모든 Python/Ansible 명령어는 `uv run` 접두사 사용
4. **보안 중심**: 민감한 정보는 Ansible Vault로 암호화

### Key Files

- `01-server-init/Justfile`: 모든 실행 명령어 정의 (그룹별 정리)
- `01-server-init/inventory/group_vars/all/vault.yml`: 암호화된 민감 변수
- `01-server-init/inventory/hosts.yml`: 호스트 정의 (IP, 사용자명만)
- `01-server-init/.vault_pass`: Vault 패스워드 파일 (gitignore 대상)
- `01-server-init/ansible.cfg`: Ansible 기본 설정

### Playbook Architecture

각 플레이북은 독립적인 변수를 가지며 특정 기능을 담당합니다:

1. **ubuntu-setup.yml**: Ubuntu 초기 설정 (APT, SSH, 타임존, LVM, Swap)
2. **wakeonlan-setup.yml**: WakeOnLAN 네트워크 설정
3. **nfs-server-setup.yml**: NFS 서버 설정 (K8s 스토리지용)
4. **k3s-setup.yml**: K3s Kubernetes 클러스터 설치
5. **slack-alert.yml**: Slack 알림 전송

### Variable Management

- **호스트 변수**: `inventory/hosts.yml`에 IP와 사용자명만 정의
- **민감 변수**: `inventory/group_vars/all/vault.yml`에 암호화 저장
- **플레이북 변수**: 각 플레이북 내부에 독립적으로 정의

## Configuration

### Vault Setup (필수)
```bash
# 1. Vault 패스워드 파일 생성
just vault-create-pass

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
- SSH 공개키 생성 (`ssh-keygen -t rsa`)
- 타겟 서버에 SSH 접근 권한

## Development Guidelines

- 모든 작업은 `01-server-init` 디렉토리에서 실행
- `just` 명령어 사용 권장 (표준화된 실행 환경)
- 새로운 플레이북 추가 시 Justfile에 `run-*`와 `check-*` 명령어 추가
- 민감한 정보는 반드시 Vault로 암호화 (`just vault-edit`)
- 플레이북별 독립적인 변수 관리 (플레이북 내부 vars 섹션 활용)