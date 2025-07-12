# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

이 프로젝트는 홈랩 환경을 자동화하는 Infrastructure as Code (IaC) 프로젝트입니다. Ansible을 기반으로 Ubuntu 서버의 초기 설정부터 K3s 클러스터 설치까지 전체 홈랩 환경을 단계별로 구성합니다.

## Project Structure

프로젝트는 단계별 구조로 구성되어 있습니다:

```
01-server-init/     # Ansible 기반 서버 초기 설정 및 K3s 설치
02-kubernetes/      # Kubernetes 리소스 정의 (향후 확장)
03-applications/    # 애플리케이션 배포 (향후 확장)
04-monitoring/      # 모니터링 스택 (향후 확장)
05-security/        # 보안 설정 (향후 확장)
```

## Development Commands

### 환경 설정
```bash
# Python 가상환경 생성 및 활성화 (uv 사용)
uv venv
source .venv/bin/activate

# 의존성 설치
uv pip install -r pyproject.toml

# 설정 파일 준비
cd 01-server-init
cp vars/main.yml.example vars/main.yml
# main.yml 파일을 실제 환경에 맞게 수정
```

### 홈랩 배포
```bash
# 전체 홈랩 환경 자동 설치 (권장)
cd 01-server-init
./run.sh

# 개별 플레이북 실행
ansible-playbook playbooks/ubuntu-setup.yml -e @vars/main.yml
ansible-playbook playbooks/wakeonlan-setup.yml -e @vars/main.yml
ansible-playbook playbooks/k3s-install.yml -e @vars/main.yml

# 연결 테스트
ansible homelab -m ping -e @vars/main.yml
```

### Ansible Commands
```bash
# 모든 작업은 01-server-init 디렉토리에서 실행
cd 01-server-init

# 인벤토리 확인
ansible-inventory --list

# 특정 호스트 정보
ansible homelab -m setup -e @vars/main.yml

# 플레이북 문법 검사
ansible-playbook --syntax-check playbooks/<playbook-name>.yml
```

## Architecture

### Core Components

1. **Ansible 기반 자동화**: Ubuntu 서버 초기 설정, WakeOnLAN 설정, K3s 클러스터 설치를 통합 관리
2. **단계별 구조**: 각 단계(01-05)는 독립적으로 관리되며 향후 확장 가능
3. **변수 기반 설정**: `vars/main.yml`에서 타겟 서버 정보와 설정을 중앙 관리

### Key Files

- `01-server-init/run.sh`: 전체 홈랩 설치를 실행하는 메인 스크립트
- `01-server-init/vars/main.yml`: 타겟 서버 정보 및 설정 변수 (main.yml.example을 복사하여 생성)
- `01-server-init/playbooks/`: Ubuntu 초기설정, WakeOnLAN, K3s 설치 플레이북
- `01-server-init/inventory/hosts.yml`: Ansible 인벤토리 설정

### Installation Flow

1. **Ubuntu 초기설정** (`ubuntu-setup.yml`):
   - APT 미러를 카카오로 변경
   - SSH 보안 설정 (공개키 인증, 패스워드 로그인 비활성화)
   - 타임존 설정
   - 필수 패키지 설치

2. **WakeOnLAN 설정** (`wakeonlan-setup.yml`):
   - 네트워크 카드 WOL 설정
   - 시스템 서비스 등록

3. **K3s 설치** (`k3s-install.yml`):
   - K3s 단일 서버 설치
   - kubeconfig 로컬 복사
   - 클러스터 상태 확인

## Configuration

### Required Variables (vars/main.yml)
```yaml
target_ip: "192.168.0.101"           # 타겟 서버 IP
target_username: "ubuntu"            # 타겟 서버 사용자명
target_password: "your_password"     # 타겟 서버 패스워드
local_public_key_path: "~/.ssh/id_rsa.pub"  # 로컬 공개키 경로
```

### Prerequisites
- Python 3.12+
- uv (Python 패키지 관리자)
- SSH 공개키 생성 (`ssh-keygen -t rsa`)
- 타겟 서버에 SSH 접근 권한

## Development Guidelines

- 모든 Ansible 작업은 `01-server-init` 디렉토리에서 실행
- 새로운 플레이북 추가 시 `run.sh`에 통합 고려
- 변수는 `vars/main.yml`에서 중앙 관리
- 민감한 정보는 절대 커밋하지 않음 (vars/main.yml은 .gitignore 대상)