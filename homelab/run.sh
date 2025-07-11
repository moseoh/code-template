#!/bin/bash

# 홈랩 전체 설치 스크립트
echo "홈랩 전체 설치를 시작합니다..."

# 가상환경 활성화
source .venv/bin/activate

# 변수 파일 확인
if [ ! -f "vars/main.yml" ]; then
    echo "오류: vars/main.yml 파일이 없습니다."
    echo "vars/main.yml.example 파일을 복사하여 vars/main.yml로 생성한 후 실제 값으로 수정하세요."
    echo "cp vars/main.yml.example vars/main.yml"
    exit 1
fi

echo "설정 확인:"
echo "- Ubuntu 24 초기설정 (APT 미러, SSH 보안, 타임존)"
echo "- K3s 클러스터 설치"
echo ""

read -p "전체 설치를 진행하시겠습니까? (y/N): " confirm
if [[ $confirm != [yY] ]]; then
    echo "설치를 취소했습니다."
    exit 1
fi

# 연결 테스트
echo "타겟 서버 연결 테스트..."
ansible homelab -m ping -e @vars/main.yml

if [ $? -ne 0 ]; then
    echo "연결 실패! 설정을 확인해주세요."
    exit 1
fi

echo "연결 성공! 설치를 시작합니다..."

# Ubuntu 초기설정
echo "===== Ubuntu 초기설정 시작 ====="
ansible-playbook playbooks/ubuntu-setup.yml -e @vars/main.yml

if [ $? -ne 0 ]; then
    echo "Ubuntu 초기설정 실패!"
    exit 1
fi

echo "===== Ubuntu 초기설정 완료 ====="
echo ""

# K3s 설치
echo "===== K3s 설치 시작 ====="
ansible-playbook playbooks/k3s-install.yml -e @vars/main.yml

if [ $? -ne 0 ]; then
    echo "K3s 설치 실패!"
    exit 1
fi

echo "===== K3s 설치 완료 ====="
echo ""
