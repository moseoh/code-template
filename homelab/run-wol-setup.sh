#!/bin/bash

# WakeOnLAN 설정 스크립트
echo "WakeOnLAN 설정을 시작합니다..."

# 가상환경 활성화
source .venv/bin/activate

# 변수 파일 확인
if [ ! -f "vars/main.yml" ]; then
    echo "오류: vars/main.yml 파일이 없습니다."
    echo "vars/main.yml.example 파일을 복사하여 vars/main.yml로 생성한 후 실제 값으로 수정하세요."
    exit 1
fi

echo "WakeOnLAN 설정 단계:"
echo "- 네트워크 인터페이스 자동 감지"
echo "- ethtool 설치 및 WOL 지원 확인"
echo "- systemd service 등록"
echo "- WOL 영구 설정 활성화"
echo ""

read -p "WakeOnLAN 설정을 진행하시겠습니까? (y/N): " confirm
if [[ $confirm != [yY] ]]; then
    echo "설정을 취소했습니다."
    exit 1
fi

# 연결 테스트
echo "타겟 서버 연결 테스트..."
ansible homelab -m ping -e @vars/main.yml

if [ $? -ne 0 ]; then
    echo "연결 실패! 설정을 확인해주세요."
    exit 1
fi

echo "연결 성공! WakeOnLAN 설정을 시작합니다..."

# WOL 설정
ansible-playbook playbooks/wakeonlan-setup.yml -e @vars/main.yml

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================="
    echo "WakeOnLAN 설정이 완료되었습니다!"
    echo "========================================="
    echo ""
    echo "사용법:"
    echo "1. 서버를 종료합니다"
    echo "2. 로컬에서 WOL 패킷 전송: wakeonlan [MAC_ADDRESS]"
    echo "3. 서버가 자동으로 부팅됩니다"
    echo ""
else
    echo "WakeOnLAN 설정 실패!"
    exit 1
fi