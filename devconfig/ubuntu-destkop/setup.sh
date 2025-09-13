#!/bin/bash

set -e

echo "================================"
echo "Ubuntu Desktop 환경 설정 시작"
echo "================================"

# 1. APT 패키지 업데이트 및 필수 패키지 설치
echo ""
echo "[1/7] APT 패키지 업데이트 및 설치..."
echo "--------------------------------"
sudo apt update

# 모든 필요한 패키지 한번에 설치
sudo apt install -y \
    zsh \
    git \
    curl \
    wget \
    fcitx5 \
    fcitx5-hangul \
    dconf-cli \
    uuid-runtime

# 2. keyd 설치 (키보드 매핑용)
echo ""
echo "[2/7] keyd 설치 및 설정..."
echo "--------------------------------"
sudo add-apt-repository -y ppa:keyd-team/ppa
sudo apt update
sudo apt install -y keyd

# keyd 설정 파일 생성 (CapsLock을 F12로 매핑)
sudo tee /etc/keyd/default.conf << 'EOF' > /dev/null
[ids]
*

[main]
capslock = f12
EOF

# keyd 서비스 활성화 및 시작
sudo systemctl enable keyd
sudo systemctl restart keyd

# 3. zsh 플러그인용 디렉토리 준비
echo ""
echo "[3/7] zsh 플러그인 다운로드..."
echo "--------------------------------"
# Oh My Zsh custom 디렉토리 미리 생성
ZSH_CUSTOM="${HOME}/.oh-my-zsh/custom"
mkdir -p "${ZSH_CUSTOM}/plugins"
mkdir -p "${ZSH_CUSTOM}/themes"

# zsh 플러그인 클론
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git "${ZSH_CUSTOM}/plugins/zsh-syntax-highlighting"
git clone https://github.com/zsh-users/zsh-autosuggestions "${ZSH_CUSTOM}/plugins/zsh-autosuggestions"

# Powerlevel10k 테마 클론
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git "${ZSH_CUSTOM}/themes/powerlevel10k"

# 4. Nerd Font 설치
echo ""
echo "[4/7] Nerd Font 설치..."
echo "--------------------------------"
mkdir -p ~/.local/share/fonts
cd ~/.local/share/fonts
curl -fLO https://github.com/ryanoasis/nerd-fonts/raw/HEAD/patched-fonts/DroidSansMono/DroidSansMNerdFont-Regular.otf
fc-cache -fv ~/.local/share/fonts/
cd -

# 5. 기본 셸을 zsh로 변경
echo ""
echo "[5/7] 기본 셸 변경..."
echo "--------------------------------"
sudo chsh -s $(which zsh) $USER

# 6. fcitx5 입력기 설정
echo ""
echo "[6/7] 한글 입력기 설정..."
echo "--------------------------------"
im-config -n fcitx5

# 7. post-reboot 스크립트 실행 권한 설정
echo ""
echo "[7/7] 후속 스크립트 준비..."
echo "--------------------------------"
chmod +x post-reboot.sh

echo ""
echo "================================"
echo "초기 설정 완료!"
echo "================================"
echo ""
echo "다음 단계:"
echo "1. 시스템을 재시작하세요: sudo reboot"
echo "2. 재시작 후 ./post-reboot.sh 를 실행하세요"
echo ""
echo "재시작 후 변경사항:"
echo "- 기본 셸이 zsh로 변경됩니다"
echo "- CapsLock 키가 F12로 매핑됩니다"
echo "- 한글 입력기(fcitx5)가 활성화됩니다"
echo ""