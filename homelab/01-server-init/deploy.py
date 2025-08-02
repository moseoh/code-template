#!/usr/bin/env python3
"""
홈랩 전체 설치 스크립트
Ubuntu 서버 초기 설정부터 K3s 클러스터 설치까지 자동화합니다.
"""

import sys
from pathlib import Path
import ansible_runner
import yaml


def check_prerequisites():
    """배포 전 사전 요구사항 확인"""
    print("📋 사전 요구사항 확인 중...")
    
    # 변수 파일 확인
    vars_file = Path(__file__).parent / "vars" / "main.yml"
    if not vars_file.exists():
        print("❌ 오류: vars/main.yml 파일이 없습니다.")
        print("   vars/main.yml.example 파일을 복사하여 vars/main.yml로 생성한 후 실제 값으로 수정하세요.")
        print("   cp vars/main.yml.example vars/main.yml")
        sys.exit(1)
    
    print("✅ vars/main.yml 파일 존재")
    
    # 플레이북 파일들 확인
    playbooks = ["ubuntu-setup.yml", "wakeonlan-setup.yml", "k3s-install.yml", "slack-alert.yml"]
    playbook_dir = Path(__file__).parent / "playbooks"
    
    for playbook in playbooks:
        playbook_path = playbook_dir / playbook
        if not playbook_path.exists():
            print(f"❌ {playbook_path} 파일이 없습니다")
            sys.exit(1)
    
    print("✅ 모든 플레이북 파일 존재")
    
    return vars_file


def test_connection(vars_file):
    """타겟 서버 연결 테스트"""
    print("\n🔍 타겟 서버 연결 테스트...")
    
    # 변수 파일 로드
    with open(vars_file, 'r') as f:
        extra_vars = yaml.safe_load(f)
    
    # ping 모듈로 연결 테스트
    r = ansible_runner.run(
        private_data_dir=str(Path(__file__).parent),
        module="ping",
        host_pattern="homelab",
        extravars=extra_vars,
        quiet=True
    )
    
    if r.status != "successful":
        print("❌ 연결 실패! 설정을 확인해주세요.")
        if r.stdout:
            print(f"출력: {r.stdout.read()}")
        sys.exit(1)
    
    print("✅ 연결 성공!")


def run_playbook(playbook_name, vars_file, description):
    """Ansible 플레이북 실행"""
    print(f"\n{'='*5} {description} 시작 {'='*5}")
    
    # 변수 파일 로드
    with open(vars_file, 'r') as f:
        extra_vars = yaml.safe_load(f)
    
    playbook_path = Path(__file__).parent / "playbooks" / playbook_name
    
    # 플레이북 실행
    r = ansible_runner.run(
        private_data_dir=str(Path(__file__).parent),
        playbook=str(playbook_path),
        extravars=extra_vars,
        verbosity=1
    )
    
    if r.status != "successful":
        print(f"❌ {description} 실패!")
        return False
    
    print(f"{'='*5} {description} 완료 {'='*5}")
    return True


def main():
    """메인 실행 함수"""
    print("🏠 홈랩 전체 설치 스크립트")
    print("=" * 50)
    
    print("\n설정 확인:")
    print("- Ubuntu 24 초기설정 (APT 미러, SSH 보안, 타임존)")
    print("- WakeOnLAN 설정")
    print("- K3s 클러스터 설치")
    print("- Slack 부팅 알림 설정 (webhook URL 설정시)")
    print("")
    
    # 사용자 확인
    confirm = input("전체 설치를 진행하시겠습니까? (y/N): ")
    if confirm.lower() not in ['y', 'yes']:
        print("설치를 취소했습니다.")
        sys.exit(1)
    
    try:
        # 사전 요구사항 확인
        vars_file = check_prerequisites()
        
        # 연결 테스트
        test_connection(vars_file)
        
        print("\n🚀 설치를 시작합니다...")
        
        # Ubuntu 초기설정
        if not run_playbook("ubuntu-setup.yml", vars_file, "Ubuntu 초기설정"):
            sys.exit(1)
        
        # WakeOnLAN 설정 (실패해도 계속 진행)
        run_playbook("wakeonlan-setup.yml", vars_file, "WakeOnLAN 설정")
        
        # K3s 설치
        if not run_playbook("k3s-install.yml", vars_file, "K3s 설치"):
            sys.exit(1)
        
        # Slack 알림 설정 (실패해도 계속 진행)
        run_playbook("slack-alert.yml", vars_file, "Slack 부팅 알림 설정")
        
        print("\n🎉 홈랩 전체 설치가 성공적으로 완료되었습니다!")
        print("\n📋 다음 단계:")
        print("1. ~/.kube/homelab-config 파일이 생성되었습니다.")
        print("2. kubectl을 사용하여 클러스터에 접근할 수 있습니다:")
        print("   export KUBECONFIG=~/.kube/homelab-config")
        print("   kubectl get nodes")
        print("3. 다음 부팅시부터 Slack 알림이 전송됩니다. (설정된 경우)")
        
    except KeyboardInterrupt:
        print("\n❌ 사용자에 의해 설치가 중단되었습니다")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ 예상치 못한 오류 발생: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()