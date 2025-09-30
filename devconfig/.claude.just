# Claude 설정 파일 관리
# devconfig/.claude와 ~/.claude 간 동기화

claude_source_dir := justfile_directory() + '/.claude'
claude_target_dir := env_var('HOME') + '/.claude'

# ~/.claude 파일들을 devconfig/.claude로 가져오기 (pull)
[group('claude')]
claude-pull:
    #!/usr/bin/env bash
    set -euo pipefail

    echo "📥 ~/.claude 파일들을 devconfig/.claude로 가져옵니다"
    echo "=================================================="

    SOURCE_DIR="{{ claude_target_dir }}"
    TARGET_DIR="{{ claude_source_dir }}"

    if [ ! -d "$SOURCE_DIR" ]; then
        echo "❌ $SOURCE_DIR 디렉토리가 존재하지 않습니다."
        exit 1
    fi

    echo ""
    echo "📁 대상 디렉토리: $TARGET_DIR"

    # 디렉토리 구조 생성
    mkdir -p "$TARGET_DIR"

    # devconfig/.claude에 존재하는 파일만 ~/.claude에서 가져오기
    echo ""
    echo "📋 파일 복사 중..."

    # devconfig/.claude에 있는 파일 목록을 기준으로 복사
    find "$TARGET_DIR" -type f | while read -r target_file; do
        # 상대 경로 계산
        rel_path="${target_file#$TARGET_DIR/}"
        source_file="$SOURCE_DIR/$rel_path"

        if [ -f "$source_file" ]; then
            # 디렉토리 생성
            mkdir -p "$(dirname "$target_file")"
            # 파일 복사
            cp -v "$source_file" "$target_file"
        else
            echo "⚠️  $source_file 파일이 존재하지 않습니다 (스킵)"
        fi
    done

    echo ""
    echo "=================================================="
    echo "✅ 가져오기 완료!"
    echo ""
    echo "복사된 위치: $TARGET_DIR"

# devconfig/.claude 파일들을 ~/.claude로 배포 (push)
[group('claude')]
claude-push:
    #!/usr/bin/env bash
    set -euo pipefail

    echo "📤 devconfig/.claude 파일들을 ~/.claude로 배포합니다"
    echo "=================================================="

    SOURCE_DIR="{{ claude_source_dir }}"
    TARGET_DIR="{{ claude_target_dir }}"

    if [ ! -d "$SOURCE_DIR" ]; then
        echo "❌ $SOURCE_DIR 디렉토리가 존재하지 않습니다."
        exit 1
    fi

    echo ""
    echo "📁 소스 디렉토리: $SOURCE_DIR"
    echo "📁 대상 디렉토리: $TARGET_DIR"

    # 대상 디렉토리 생성
    mkdir -p "$TARGET_DIR"

    # devconfig/.claude에 있는 파일만 ~/.claude로 복사 (추가 파일은 유지)
    echo ""
    echo "📋 파일 복사 중..."

    # devconfig/.claude의 모든 파일을 ~/.claude로 복사
    find "$SOURCE_DIR" -type f | while read -r source_file; do
        # 상대 경로 계산
        rel_path="${source_file#$SOURCE_DIR/}"
        target_file="$TARGET_DIR/$rel_path"

        # 디렉토리 생성
        mkdir -p "$(dirname "$target_file")"
        # 파일 복사
        cp -v "$source_file" "$target_file"
    done

    echo ""
    echo "=================================================="
    echo "✅ 배포 완료!"
    echo ""
    echo "배포된 위치: $TARGET_DIR"

# devconfig/.claude 파일 목록 확인
[group('claude')]
claude-list:
    #!/usr/bin/env bash
    set -euo pipefail

    echo "📋 devconfig/.claude 파일 목록"
    echo "=================================================="

    SOURCE_DIR="{{ claude_source_dir }}"

    if [ ! -d "$SOURCE_DIR" ]; then
        echo "❌ $SOURCE_DIR 디렉토리가 존재하지 않습니다."
        exit 1
    fi

    echo ""
    tree "$SOURCE_DIR" 2>/dev/null || find "$SOURCE_DIR" -type f

    echo ""
    echo "=================================================="