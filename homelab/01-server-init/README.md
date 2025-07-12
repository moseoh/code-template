# 호스트 서버 설정

홈랩 서버 초기 설정 및 K3s 클러스터 자동 설치

## 기술 스택

- **Ansible**: 서버 자동화 및 구성 관리
- **Python 3.12+**: 런타임 환경
- **uv**: Python 패키지 관리

## 기능

- Ubuntu 서버 초기 설정 (SSH, 보안, 패키지)
- WakeOnLAN 네트워크 설정
- K3s Kubernetes 클러스터 설치
- 설정 중앙화 및 자동화

## 사전 준비

```bash
# Python 3.12+ 및 uv 설치 필요
python3 --version  # 3.12+
uv --version
```

## 빠른 시작

```bash
./run.sh
```
