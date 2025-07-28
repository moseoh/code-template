#!/usr/bin/env python3
"""
GitHub Repository Ruleset 설정 스크립트

이 스크립트는 GitHub CLI를 사용하여 repository에 기본적인 ruleset을 설정합니다.
- PR 필수 설정
- Squash merge 설정
- 기본적인 브랜치 보호 규칙 적용
"""

import json
import subprocess
import sys
import os

def run_command(cmd, capture_output=True, check=True):
    """명령어 실행 헬퍼 함수"""
    try:
        result = subprocess.run(
            cmd, 
            shell=True, 
            capture_output=capture_output, 
            text=True, 
            check=check
        )
        return result
    except subprocess.CalledProcessError as e:
        print(f"❌ 명령어 실행 실패: {cmd}")
        print(f"오류: {e.stderr if e.stderr else e}")
        sys.exit(1)

def get_repo_info():
    """현재 디렉토리의 GitHub repository 정보 및 기본 브랜치 획득"""
    print("\n🔍 Repository 정보 확인 중...")
    
    # git remote -v로 repository 정보 확인
    result = run_command("git remote -v")
    if not result.stdout:
        print("❌ Git repository가 아니거나 remote가 설정되지 않았습니다.")
        sys.exit(1)
    
    # GitHub repository 정보 추출
    for line in result.stdout.strip().split('\n'):
        if 'github.com' in line and 'origin' in line and '(fetch)' in line:
            # git@github.com:owner/repo.git 또는 https://github.com/owner/repo.git 형태 파싱
            url = line.split()[1]
            if url.startswith('git@github.com:'):
                repo_path = url.replace('git@github.com:', '').replace('.git', '')
            elif url.startswith('https://github.com/'):
                repo_path = url.replace('https://github.com/', '').replace('.git', '')
            else:
                continue

            owner, repo = repo_path.split('/')
            print(f"📁 Repository: {owner}/{repo}")
            return owner, repo
    
    print("❌ GitHub repository를 찾을 수 없습니다.")
    sys.exit(1)

def get_default_branch(owner, repo):
    """GitHub API를 통해 기본 브랜치 확인"""
    print("\n🔍 기본 브랜치 확인 중...")
    
    try:
        result = run_command(f"gh api repos/{owner}/{repo}")
        repo_info = json.loads(result.stdout)
        default_branch = repo_info.get('default_branch', 'main')
        print(f"📌 기본 브랜치: {default_branch}")
        return default_branch
    except Exception as e:
        print(f"⚠️  기본 브랜치 확인 실패, main으로 설정: {e}")
        return 'main'

def check_existing_rulesets(owner, repo):
    """기존 ruleset 존재 여부 확인 및 ID 반환"""
    print("\n🔍 기존 ruleset 확인 중...")
    
    try:
        result = run_command(f"gh api repos/{owner}/{repo}/rulesets", check=False)
        if result.returncode == 0:
            rulesets = json.loads(result.stdout)
            if rulesets and len(rulesets) > 0:
                print(f"✅ 이미 {len(rulesets)}개의 ruleset이 존재합니다:")
                for ruleset in rulesets:
                    print(f"  - {ruleset.get('name', 'Unnamed')} (ID: {ruleset.get('id')})")
                    if ruleset.get('name') == 'Default Branch Protection':
                        print(f"🔄 'Default Branch Protection' ruleset을 업데이트합니다.")
                        return ruleset.get('id')
                return True  # 다른 ruleset이 있으면 종료
        return False
    except json.JSONDecodeError:
        print("⚠️  ruleset 정보를 파싱할 수 없습니다. 계속 진행합니다.")
        return False
    except Exception as e:
        print(f"⚠️  ruleset 확인 중 오류 발생: {e}. 계속 진행합니다.")
        return False

def create_or_update_ruleset(owner, repo, ruleset_id=None):
    """기본 ruleset 생성 또는 업데이트"""
    action = "업데이트" if ruleset_id else "생성"
    print(f"\n📝 기본 ruleset {action} 중...")
    
    # 기본 ruleset 설정
    ruleset_config = {
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
                    "dismiss_stale_reviews_on_push": True,
                    "require_code_owner_review": False,
                    "require_last_push_approval": False,
                    "required_approving_review_count": 1,
                    "required_review_thread_resolution": False,
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
                "actor_id": 5 # repository admin
            }
        ]
    }
    
    # JSON 파일로 임시 저장
    config_file = "/tmp/github_ruleset.json"
    squash_file = "/tmp/github_squash.json"  # 변수를 미리 선언
    
    with open(config_file, 'w', encoding='utf-8') as f:
        json.dump(ruleset_config, f, indent=2)
    
    try:
        # ruleset 생성 또는 업데이트
        if ruleset_id:
            result = run_command(f"gh api repos/{owner}/{repo}/rulesets/{ruleset_id} --method PUT --input {config_file}")
            print("✅ Ruleset이 성공적으로 업데이트되었습니다!")
        else:
            result = run_command(f"gh api repos/{owner}/{repo}/rulesets --method POST --input {config_file}")
            print("✅ Ruleset이 성공적으로 생성되었습니다!")
        
        # Repository 설정에서 squash merge 활성화
        print("\n📝 Squash merge 설정 중...")
        squash_config = {
            "allow_squash_merge": True,
            "allow_merge_commit": False,
            "allow_rebase_merge": False,
            "delete_branch_on_merge": True,
            "squash_merge_commit_title": "PR_TITLE",
            "squash_merge_commit_message": "PR_BODY"
        }
        
        # squash_file은 이미 위에서 선언됨
        with open(squash_file, 'w', encoding='utf-8') as f:
            json.dump(squash_config, f, indent=2)
        
        run_command(f"gh api repos/{owner}/{repo} --method PATCH --input {squash_file}")
        print("✅ Squash merge가 활성화되었습니다!")
        
    except Exception as e:
        print(f"❌ Ruleset {action} 실패: {e}")
        sys.exit(1)
    finally:
        # 임시 파일 정리
        for temp_file in [config_file, squash_file]:
            if os.path.exists(temp_file):
                os.remove(temp_file)


def main():
    """메인 함수"""
    print("🚀 GitHub Repository Ruleset 설정 스크립트")
    print("=" * 50)
    
    # GitHub CLI 설치 확인
    print("\n🔍 GitHub CLI 설치 상태 확인 중...")
    result = run_command("which gh", check=False)
    if result.returncode != 0:
        print("❌ GitHub CLI (gh)가 설치되지 않았습니다.")
        print()
        print("📦 설치 방법:")
        print("  brew install gh")
        print()
        print("설치 후 다음 명령어로 로그인하세요:")
        print("  gh auth login")
        print()
        print("자세한 정보: https://cli.github.com/")
        sys.exit(1)
    else:
        print("✅ GitHub CLI가 설치되어 있습니다.")
    
    # GitHub CLI 인증 확인
    print("\n🔍 GitHub CLI 인증 상태 확인 중...")
    result = run_command("gh auth status", check=False)
    # gh auth status는 성공 시에도 stderr로 출력하므로 실제 API 호출로 확인
    test_result = run_command("gh api user", check=False)
    if test_result.returncode != 0:
        print("❌ GitHub CLI에 로그인되지 않았습니다.")
        print("로그인 방법: gh auth login")
        sys.exit(1)
    else:
        print("✅ GitHub CLI가 인증되었습니다.")
    
    # Repository 정보 확인
    owner, repo = get_repo_info()
    
    # 기본 브랜치 확인
    default_branch = get_default_branch(owner, repo)
    
    # 기존 ruleset 확인
    existing_ruleset = check_existing_rulesets(owner, repo)
    if existing_ruleset is True:  # 다른 ruleset이 있으면 종료
        return
    
    # Ruleset 생성 또는 업데이트
    ruleset_id = existing_ruleset if isinstance(existing_ruleset, int) else None
    create_or_update_ruleset(owner, repo, ruleset_id)
    
    print("\n" + "=" * 50)
    print("🎉 GitHub repository ruleset 설정이 완료되었습니다!")
    print(f"🌐 확인: https://github.com/{owner}/{repo}/settings/rules")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n⚠️  사용자에 의해 중단되었습니다.")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ 예상치 못한 오류 발생: {e}")
        sys.exit(1)