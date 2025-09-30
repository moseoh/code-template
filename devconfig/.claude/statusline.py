#!/usr/bin/env python3
import json
import sys
import subprocess
import os
from datetime import datetime

def get_git_branch():
    """현재 디렉토리의 git 브랜치를 가져옵니다."""
    try:
        result = subprocess.run(
            ['git', 'branch', '--show-current'], 
            capture_output=True, 
            text=True, 
            timeout=2
        )
        if result.returncode == 0 and result.stdout.strip():
            return result.stdout.strip()
        return "no-git"
    except (subprocess.TimeoutExpired, FileNotFoundError):
        return "no-git"

def get_project_name():
    """현재 작업 디렉토리의 프로젝트 이름을 가져옵니다."""
    return os.path.basename(os.getcwd())

def main():
    try:
        # stdin에서 JSON 데이터 읽기
        input_data = json.loads(sys.stdin.read())
        
        # 현재 시간 (한국 형식)
        current_time = datetime.now().strftime("%H:%M:%S")
        
        # 모델 이름
        model_name = input_data.get('model', {}).get('display_name', 'Claude')
        
        # 프로젝트 이름
        project_name = get_project_name()
        
        # Git 브랜치
        git_branch = get_git_branch()
        
        # 상태라인 구성 (이모지와 파이프 구분자 사용)
        status_parts = [
            f"⏰ {current_time}",
            f"🤖 {model_name}",
            f"📁 {project_name}",
            f"🌿 {git_branch}"
        ]
        
        # 파이프로 구분하여 출력
        status_line = " | ".join(status_parts)
        print(f"\033[36m{status_line}\033[0m")  # 청록색으로 출력
        
    except Exception as e:
        # 오류 발생 시 기본 상태라인 출력
        print("🤖 Claude Code | 📁 Project")

if __name__ == "__main__":
    main()