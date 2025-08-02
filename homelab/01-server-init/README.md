# 서버 초기화

홈랩 Ubuntu 서버 초기 설정 및 K3s 클러스터 자동 설치

## 개요

Ansible을 사용하여 Ubuntu 서버의 초기 설정부터 K3s Kubernetes 클러스터 설치까지 완전 자동화합니다.

## 기술 스택

- **Ansible**: 서버 자동화 및 구성 관리
- **ansible-runner**: Python 기반 Ansible 실행
- **Python 3.12+**: 런타임 환경
- **uv**: Python 패키지 관리

## 주요 기능

- Ubuntu 서버 초기 설정 (APT 미러, SSH 보안, 타임존)
- WakeOnLAN 네트워크 설정
- K3s Kubernetes 클러스터 설치
- kubeconfig 자동 설정
- 설정 중앙화 및 자동화

## 사전 준비

### 필수 소프트웨어
```bash
# Python 3.12+ 및 uv 설치 확인
python3 --version  # 3.12+
uv --version

# SSH 키 생성 (없는 경우)
ssh-keygen -t rsa
```

### 설정 파일 준비
```bash
# 변수 파일 복사 및 수정
cp vars/main.yml.example vars/main.yml
# vars/main.yml 파일을 실제 환경에 맞게 수정
```

## 실행 방법

### 전체 설치 (권장)
```bash
# 의존성 설치 및 전체 배포
uv run deploy.py
```

### 단계별 실행
```bash
# 개별 플레이북 실행
ansible-playbook playbooks/ubuntu-setup.yml -e @vars/main.yml
ansible-playbook playbooks/wakeonlan-setup.yml -e @vars/main.yml
ansible-playbook playbooks/k3s-install.yml -e @vars/main.yml
```

### 연결 테스트
```bash
# 타겟 서버 연결 확인
ansible homelab -m ping -e @vars/main.yml
```

## 설치 완료 후

```bash
# kubeconfig 설정
export KUBECONFIG=~/.kube/homelab-config

# 클러스터 상태 확인
kubectl get nodes
kubectl get pods --all-namespaces
```
