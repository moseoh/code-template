#!/usr/bin/env python3
"""
Claude 설정 백업 스크립트
~/.claude 폴더의 설정 파일들을 현재 프로젝트의 claude 폴더로 백업합니다.
"""

import os
import shutil
import fnmatch
from pathlib import Path


def get_excluded_patterns():
    """제외할 파일/폴더 패턴 목록"""
    return [
        '.DS_Store',
        '*.tmp',
        '*.temp',
        '*.cache',
        '*.log',
        '__pycache__',
        '.vscode',
        'ide',          # IDE 관련 설정
        'projects',     # 프로젝트별 설정 
        'shell-snapshots',  # 개인 쉘 스냅샷
        'todos',        # 개인 todo 목록
        'statsig'       # 개인 통계
    ]


def should_exclude(path, patterns):
    """주어진 경로가 제외 패턴에 매칭되는지 확인"""
    path_name = os.path.basename(path)
    for pattern in patterns:
        if fnmatch.fnmatch(path_name, pattern):
            return True
    return False


def backup_claude_config():
    """~/.claude 설정을 현재 프로젝트로 백업"""
    source_dir = Path.home() / '.claude'
    target_dir = Path(__file__).parent / '.claude'
    
    if not source_dir.exists():
        print(f"❌ 소스 디렉토리가 존재하지 않습니다: {source_dir}")
        return False
    
    # 타겟 디렉토리 생성
    target_dir.mkdir(exist_ok=True)
    
    excluded_patterns = get_excluded_patterns()
    copied_files = []
    skipped_files = []
    
    print(f"📁 백업 시작: {source_dir} -> {target_dir}")
    print(f"🚫 제외 패턴: {', '.join(excluded_patterns)}")
    print()
    
    # 재귀적으로 파일 복사
    for root, dirs, files in os.walk(source_dir):
        # 제외할 디렉토리 필터링
        dirs[:] = [d for d in dirs if not should_exclude(os.path.join(root, d), excluded_patterns)]
        
        for file in files:
            source_file = os.path.join(root, file)
            
            # 제외 패턴 확인
            if should_exclude(source_file, excluded_patterns):
                skipped_files.append(os.path.relpath(source_file, source_dir))
                continue
            
            # 상대 경로 계산
            rel_path = os.path.relpath(source_file, source_dir)
            target_file = target_dir / rel_path
            
            # 타겟 디렉토리 생성
            target_file.parent.mkdir(parents=True, exist_ok=True)
            
            # 파일 복사
            shutil.copy2(source_file, target_file)
            copied_files.append(rel_path)
    
    # 결과 출력
    print(f"✅ 복사된 파일: {len(copied_files)}개")
    for file in sorted(copied_files):
        print(f"   📄 {file}")
    
    if skipped_files:
        print(f"\n⏭️  제외된 파일: {len(skipped_files)}개")
        for file in sorted(skipped_files):
            print(f"   🚫 {file}")
    
    print(f"\n🎉 백업 완료!")
    return True


if __name__ == '__main__':
    backup_claude_config()