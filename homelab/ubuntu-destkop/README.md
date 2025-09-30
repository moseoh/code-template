# Ubuntu Desktop 개발 환경 자동 설정

Ubuntu Desktop 환경을 위한 자동화된 설정 스크립트입니다.

## 🚀 빠른 시작

```bash
# 1. 스크립트 실행 권한 부여
chmod +x setup.sh post-reboot.sh

# 2. 초기 설정 실행
./setup.sh

# 3. 시스템 재시작
sudo reboot

# 4. 재시작 후 설정 완료
./post-reboot.sh
```

## 📦 설치되는 구성 요소

### 1단계: setup.sh (재시작 전)
- **Zsh & Oh My Zsh**: 강력한 셸 환경
- **플러그인**: 
  - zsh-syntax-highlighting (명령어 하이라이팅)
  - zsh-autosuggestions (자동 완성)
- **Powerlevel10k**: 모던한 터미널 프롬프트 테마
- **Nerd Font**: 아이콘 지원 폰트 (DroidSansMono)
- **fcitx5**: 한글 입력기
- **keyd**: 키보드 매핑 (CapsLock → F12)

### 2단계: post-reboot.sh (재시작 후)
- Oh My Zsh 설치 완료
- .zshrc 자동 설정
- 터미널 테마 설치 (Gogh)
- Powerlevel10k 초기 설정

## 🔧 주요 기능

### 한글 입력
- **전환**: `Shift+Space` 또는 `한/영` 키
- **설정**: `fcitx5-configtool` 명령으로 상세 설정

### 키보드 매핑
- **CapsLock → F12**: 개발 시 유용한 단축키로 활용
- 설정 파일: `/etc/keyd/default.conf`

### 터미널 개선
- 명령어 하이라이팅
- 자동 완성 제안
- Git 상태 표시
- 아이콘 지원

## 📝 수동 설정 (선택사항)

### 터미널 폰트 변경
1. 터미널 설정 열기
2. 폰트를 "DroidSansMono Nerd Font"로 변경
3. 아이콘이 제대로 표시되는지 확인

### Powerlevel10k 재설정
```bash
p10k configure
```

### 추가 zsh 플러그인 설치
```bash
# ~/.zshrc 파일의 plugins 섹션에 추가
plugins=(
    git
    zsh-syntax-highlighting
    zsh-autosuggestions
    # 원하는 플러그인 추가
)
```

## 🗂 파일 구조

```
.
├── setup.sh          # 메인 설치 스크립트
├── post-reboot.sh    # 재시작 후 실행 스크립트
├── README.md         # 이 문서
└── old/              # 이전 개별 스크립트들 (참고용)
    ├── 1.zsh-install.sh
    ├── 2.oh-my-zsh-install.sh
    ├── 3.zsh-plugin-install/
    └── 4.hangul/
```

## ⚠️ 주의사항

- 스크립트는 Ubuntu 22.04 이상에서 테스트되었습니다
- 재시작은 **1회만** 필요합니다
- 기존 `.zshrc` 파일은 자동으로 백업됩니다
- sudo 권한이 필요합니다

## 🐛 문제 해결

### zsh가 기본 셸로 설정되지 않는 경우
```bash
chsh -s $(which zsh)
# 로그아웃 후 다시 로그인
```

### 한글 입력이 작동하지 않는 경우
```bash
im-config -n fcitx5
# 로그아웃 후 다시 로그인
```

### 플러그인이 로드되지 않는 경우
```bash
# 플러그인 경로 확인
ls ~/.oh-my-zsh/custom/plugins/
# .zshrc 파일 확인
cat ~/.zshrc | grep plugins
```

## 📄 라이선스

이 프로젝트는 개인 사용을 위한 설정 스크립트입니다.