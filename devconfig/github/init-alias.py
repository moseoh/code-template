#!/usr/bin/env python3
"""
Alias 초기화 스크립트

이 스크립트는 ~/.zshrc에 개발용 alias를 등록합니다.
"""

import os
import sys

def add_alias_to_zshrc():
    """~/.zshrc에 alias 추가"""
    zshrc_path = os.path.expanduser("~/.zshrc")
    script_dir = os.path.dirname(os.path.abspath(__file__))
    
    aliases = [
        f"alias initgithubruleset='{script_dir}/init-github-ruleset.py'"
    ]
    
    print("🔧 ~/.zshrc에 alias 추가 중...")
    
    # 기존 파일 읽기
    existing_content = ""
    if os.path.exists(zshrc_path):
        with open(zshrc_path, 'r', encoding='utf-8') as f:
            existing_content = f.read()
    
    # 새로운 alias 추가
    new_aliases = []
    for alias in aliases:
        if alias not in existing_content:
            new_aliases.append(alias)
    
    if new_aliases:
        with open(zshrc_path, 'a', encoding='utf-8') as f:
            f.write("\n# GitHub DevConfig Aliases\n")
            for alias in new_aliases:
                f.write(f"{alias}\n")
        
        print(f"✅ {len(new_aliases)}개의 새로운 alias가 추가되었습니다:")
        for alias in new_aliases:
            print(f"  - {alias}")
        
        print(f"\n🚀 설정 적용을 위해 다음 명령어를 실행하세요:")
        print(f"source ~/.zshrc")
    else:
        print("✅ 모든 alias가 이미 설정되어 있습니다.")

def main():
    """메인 함수"""
    print("🚀 GitHub DevConfig Alias 초기화")
    print("=" * 40)
    
    add_alias_to_zshrc()

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n⚠️  사용자에 의해 중단되었습니다.")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ 예상치 못한 오류 발생: {e}")
        sys.exit(1)