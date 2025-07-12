#!/usr/bin/env python3
"""
Claude 설정 업데이트 스크립트
현재 프로젝트의 claude 폴더에서 ~/.claude 폴더로 설정을 업데이트합니다.
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
        '*.py',         # 이 스크립트들은 제외
        '*.pyc'
    ]


def should_exclude(path, patterns):
    """주어진 경로가 제외 패턴에 매칭되는지 확인"""
    path_name = os.path.basename(path)
    for pattern in patterns:
        if fnmatch.fnmatch(path_name, pattern):
            return True
    return False


def update_claude_config():
    """현재 프로젝트에서 ~/.claude 설정으로 업데이트"""
    source_dir = Path(__file__).parent / '.claude'
    target_dir = Path.home() / '.claude'
    
    # 백업할 파일들만 포함하는 폴더들 확인
    config_files = list(source_dir.glob('**/*'))
    config_files = [f for f in config_files if f.is_file()]
    
    if not config_files:
        print(f"❌ 업데이트할 설정 파일이 없습니다: {source_dir}")
        return False
    
    # 타겟 디렉토리 생성
    target_dir.mkdir(exist_ok=True)
    
    excluded_patterns = get_excluded_patterns()
    updated_files = []
    skipped_files = []
    
    print(f"📁 업데이트 시작: {source_dir} -> {target_dir}")
    print(f"🚫 제외 패턴: {', '.join(excluded_patterns)}")
    print()
    
    # 재귀적으로 파일 복사
    for root, dirs, files in os.walk(source_dir):
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
            
            # 파일이 이미 존재하는지 확인
            action = "업데이트"
            if target_file.exists():
                # 파일 내용 비교
                with open(source_file, 'rb') as sf, open(target_file, 'rb') as tf:
                    if sf.read() == tf.read():
                        action = "동일함"
                    else:
                        action = "덮어씀"
            else:
                action = "새로생성"
            
            # 파일 복사
            shutil.copy2(source_file, target_file)
            updated_files.append((rel_path, action))
    
    # 결과 출력
    print(f"✅ 처리된 파일: {len(updated_files)}개")
    for file, action in sorted(updated_files):
        emoji = {
            "새로생성": "🆕",
            "덮어씀": "🔄", 
            "업데이트": "📝",
            "동일함": "✨"
        }.get(action, "📄")
        print(f"   {emoji} {file} ({action})")
    
    if skipped_files:
        print(f"\n⏭️  제외된 파일: {len(skipped_files)}개")
        for file in sorted(skipped_files):
            print(f"   🚫 {file}")
    
    print(f"\n🎉 업데이트 완료!")
    return True


if __name__ == '__main__':
    update_claude_config()