
# GitHub CLI 설치 및 인증 확인
[group('git')]
check-gh:
    #!/usr/bin/env bash
    set -euo pipefail

    echo "🔍 GitHub CLI 설치 상태 확인 중..."
    if ! command -v gh &> /dev/null; then
        echo "❌ GitHub CLI (gh)가 설치되지 않았습니다."
        echo ""
        echo "📦 설치 방법:"
        echo "  brew install gh"
        echo ""
        echo "설치 후 다음 명령어로 로그인하세요:"
        echo "  gh auth login"
        echo ""
        echo "자세한 정보: https://cli.github.com/"
        exit 1
    fi
    echo "✅ GitHub CLI가 설치되어 있습니다."

    echo ""
    echo "🔍 GitHub CLI 인증 상태 확인 중..."
    if ! gh api user &> /dev/null; then
        echo "❌ GitHub CLI에 로그인되지 않았습니다."
        echo "로그인 방법: gh auth login"
        exit 1
    fi
    echo "✅ GitHub CLI가 인증되었습니다."

# 기존 ruleset 확인
[group('git')]
check-rulesets:
    #!/usr/bin/env bash
    set -euo pipefail

    REMOTE_URL=$(git remote get-url origin)
    if [[ $REMOTE_URL == git@github.com:* ]]; then
        REPO_PATH="${REMOTE_URL#git@github.com:}"
        REPO_PATH="${REPO_PATH%.git}"
    else
        REPO_PATH="${REMOTE_URL#https://github.com/}"
        REPO_PATH="${REPO_PATH%.git}"
    fi

    echo ""
    echo "🔍 기존 ruleset 확인 중..."

    RULESETS=$(gh api "repos/$REPO_PATH/rulesets" 2>/dev/null || echo "[]")

    if [ "$RULESETS" != "[]" ]; then
        echo "✅ 기존 ruleset이 존재합니다:"
        echo "$RULESETS" | jq -r '.[] | "  - \(.name) (ID: \(.id))"'
    else
        echo "ℹ️  기존 ruleset이 없습니다."
    fi

# GitHub repository ruleset 설정 (생성 또는 업데이트)
[group('git')]
update-ruleset:
    #!/usr/bin/env bash
    set -euo pipefail

    echo "🚀 GitHub Repository Ruleset 설정 스크립트"
    echo "=================================================="

    # Repository 정보 추출
    REMOTE_URL=$(git remote get-url origin)
    if [[ $REMOTE_URL == git@github.com:* ]]; then
        REPO_PATH="${REMOTE_URL#git@github.com:}"
        REPO_PATH="${REPO_PATH%.git}"
    else
        REPO_PATH="${REMOTE_URL#https://github.com/}"
        REPO_PATH="${REPO_PATH%.git}"
    fi

    echo ""
    echo "📁 Repository: $REPO_PATH"

    # 기존 ruleset 확인
    echo ""
    echo "🔍 기존 ruleset 확인 중..."
    RULESETS=$(gh api "repos/$REPO_PATH/rulesets" 2>/dev/null || echo "[]")
    RULESET_ID=$(echo "$RULESETS" | jq -r '.[] | select(.name == "Default Branch Protection") | .id' | head -n1)

    if [ -n "$RULESET_ID" ] && [ "$RULESET_ID" != "null" ]; then
        echo "🔄 'Default Branch Protection' ruleset을 업데이트합니다. (ID: $RULESET_ID)"
        ACTION="업데이트"
        METHOD="PUT"
        ENDPOINT="repos/$REPO_PATH/rulesets/$RULESET_ID"
    else
        echo "📝 새로운 ruleset을 생성합니다."
        ACTION="생성"
        METHOD="POST"
        ENDPOINT="repos/$REPO_PATH/rulesets"
    fi

    # Ruleset 설정 JSON 생성
    cat > /tmp/github_ruleset.json << 'RULESET_EOF'
    {
      "name": "Default Branch Protection",
      "target": "branch",
      "enforcement": "active",
      "conditions": {
        "ref_name": {
          "include": ["~DEFAULT_BRANCH"],
          "exclude": []
        }
      },
      "rules": [
        {
          "type": "pull_request",
          "parameters": {
            "dismiss_stale_reviews_on_push": true,
            "require_code_owner_review": false,
            "require_last_push_approval": false,
            "required_approving_review_count": 1,
            "required_review_thread_resolution": false,
            "allowed_merge_methods": ["squash"]
          }
        },
        {
          "type": "non_fast_forward",
          "parameters": {}
        },
        {
          "type": "deletion",
          "parameters": {}
        }
      ],
      "bypass_actors": [
        {
          "actor_type": "RepositoryRole",
          "bypass_mode": "pull_request",
          "actor_id": 5
        }
      ]
    }
    RULESET_EOF

    # Ruleset 생성 또는 업데이트
    echo ""
    echo "📝 기본 ruleset ${ACTION} 중..."
    gh api "$ENDPOINT" --method "$METHOD" --input /tmp/github_ruleset.json
    echo "✅ Ruleset이 성공적으로 ${ACTION}되었습니다!"

    # Squash merge 설정
    echo ""
    echo "📝 Squash merge 설정 중..."
    cat > /tmp/github_squash.json << 'SQUASH_EOF'
    {
      "allow_squash_merge": true,
      "allow_merge_commit": false,
      "allow_rebase_merge": false,
      "delete_branch_on_merge": true,
      "squash_merge_commit_title": "PR_TITLE",
      "squash_merge_commit_message": "PR_BODY"
    }
    SQUASH_EOF

    gh api "repos/$REPO_PATH" --method PATCH --input /tmp/github_squash.json
    echo "✅ Squash merge가 활성화되었습니다!"

    # 임시 파일 정리
    rm -f /tmp/github_ruleset.json /tmp/github_squash.json

    echo ""
    echo "=================================================="
    echo "🎉 GitHub repository ruleset 설정이 완료되었습니다!"
    echo "🌐 확인: https://github.com/$REPO_PATH/settings/rules"

# Alias clean-branches
[group('git')]
cleanb: clean-branches

# 브랜치 정리 (main, master, develop, staging, release 제외)
[group('git')]
clean-branches:
    #!/usr/bin/env bash
    set -euo pipefail

    # 보호할 브랜치 목록
    PROTECTED_BRANCHES="main|master|develop|staging|release"

    echo "🧹 브랜치 정리 중..."
    echo "=================================================="
    echo ""

    git branch | grep -Ev "(^\*|$PROTECTED_BRANCHES)" | xargs --no-run-if-empty git branch -D

    echo ""
    echo "=================================================="
    echo "🎉 브랜치 정리 완료!"

# Alias quick-push
[group('git')]
qp: quick-push

# 빠른 push
[group('git')]
quick-push:
    #!/usr/bin/env bash
    set -euo pipefail

    echo "🧹 빠른 푸쉬 중..."
    echo "=================================================="
    echo ""

    git pull
    git add .
    git commit -m "update"
    git push

    echo ""
    echo "=================================================="
    echo "🎉 빠른 푸쉬 완료!"