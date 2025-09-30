#!/bin/bash

set -e

echo "================================"
echo "재시작 후 설정 스크립트"
echo "================================"

# 1. Oh My Zsh 설치
echo ""
echo "[1/4] Oh My Zsh 설치..."
echo "--------------------------------"
# Oh My Zsh 설치 (unattended mode)
if [ ! -d "$HOME/.oh-my-zsh" ]; then
    RUNZSH=no sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended
else
    echo "Oh My Zsh가 이미 설치되어 있습니다."
fi

# 2. .zshrc 파일 자동 수정
echo ""
echo "[2/4] .zshrc 파일 설정..."
echo "--------------------------------"

# 백업 생성
cp ~/.zshrc ~/.zshrc.backup.$(date +%Y%m%d_%H%M%S)

# ZSH_THEME 변경
sed -i 's/^ZSH_THEME=.*/ZSH_THEME="powerlevel10k\/powerlevel10k"/' ~/.zshrc

# plugins 설정 수정
# 먼저 기존 plugins 라인을 찾아서 주석 처리
sed -i '/^plugins=(/,/^)/ s/^/#/' ~/.zshrc

# 새로운 plugins 설정 추가
cat >> ~/.zshrc << 'EOF'

# Custom plugins configuration
plugins=(
    git
    zsh-syntax-highlighting
    zsh-autosuggestions
)

# Source oh-my-zsh
source $ZSH/oh-my-zsh.sh
EOF

# 3. 터미널 테마 설치
echo ""
echo "[3/4] 터미널 테마 설치..."
echo "--------------------------------"
echo ""
echo "Gogh 터미널 테마 설치를 시작합니다."
echo "추천 테마: Snazzy (번호: 291)"
echo ""
echo "테마 선택 프롬프트가 나타나면 원하는 테마 번호를 입력하세요."
echo "설치를 건너뛰려면 Ctrl+C를 누르세요."
echo ""
read -p "계속하려면 Enter를 누르세요..." -t 10 || true

# Gogh 터미널 테마 설치 (선택적)
if command -v gnome-terminal &> /dev/null; then
    bash -c "$(wget -qO- https://git.io/vQgMr)" || true
else
    echo "GNOME 터미널이 감지되지 않았습니다. 테마 설치를 건너뜁니다."
fi

# 4. Powerlevel10k 설정
echo ""
echo "[4/4] Powerlevel10k 초기 설정..."
echo "--------------------------------"

# p10k 설정 파일이 없으면 기본 설정 생성
if [ ! -f ~/.p10k.zsh ]; then
    echo ""
    echo "Powerlevel10k 설정 마법사가 다음 zsh 실행 시 자동으로 시작됩니다."
    echo "원하는 스타일을 선택하여 설정을 완료하세요."
    
    # p10k instant prompt 활성화
    cat >> ~/.zshrc << 'EOF'

# Enable Powerlevel10k instant prompt
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh
EOF
fi

echo ""
echo "================================"
echo "설정 완료!"
echo "================================"
echo ""
echo "다음 단계:"
echo "1. 새 터미널을 열거나 'zsh' 명령을 실행하세요"
echo "2. Powerlevel10k 설정 마법사가 자동으로 시작됩니다"
echo "3. 원하는 프롬프트 스타일을 선택하세요"
echo ""
echo "추가 설정:"
echo "- Powerlevel10k 재설정: p10k configure"
echo "- 터미널 폰트를 'DroidSansMono Nerd Font'로 변경하면 아이콘이 제대로 표시됩니다"
echo ""
echo "한글 입력:"
echo "- Shift+Space 또는 한/영 키로 전환"
echo "- fcitx5-configtool로 상세 설정 가능"
echo ""