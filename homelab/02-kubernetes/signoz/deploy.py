#!/usr/bin/env python3
"""
SigNoz 배포 스크립트
홈 서버 K3s 클러스터에 SigNoz 모니터링 스택을 배포합니다.
"""

import subprocess
import sys
from pathlib import Path

# 글로벌 설정
NAMESPACE = "signoz"


def run_command(cmd: list[str], check: bool = True) -> subprocess.CompletedProcess:
    """명령어 실행 및 결과 반환"""
    print(f"실행 중: {' '.join(cmd)}")
    result = subprocess.run(cmd, capture_output=True, text=True, check=check)
    
    if result.stdout:
        print(f"출력: {result.stdout}")
    if result.stderr:
        print(f"에러: {result.stderr}")
    
    return result


def check_prerequisites():
    """배포 전 사전 요구사항 확인"""
    print("📋 사전 요구사항 확인 중...")
    
    # kubectl 설치 확인
    try:
        run_command(["kubectl", "version", "--client"])
        print("✅ kubectl 설치됨")
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("❌ kubectl이 설치되지 않았습니다")
        sys.exit(1)
    
    # helm 설치 확인
    try:
        run_command(["helm", "version", "--short"])
        print("✅ helm 설치됨")
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("❌ helm이 설치되지 않았습니다")
        sys.exit(1)
    
    # K3s 클러스터 연결 확인
    try:
        run_command(["kubectl", "get", "nodes"])
        print("✅ K3s 클러스터 연결됨")
    except subprocess.CalledProcessError:
        print("❌ K3s 클러스터에 연결할 수 없습니다")
        sys.exit(1)
    
    # values 파일들 확인
    signoz_values = Path(__file__).parent / "signoz-values.yaml"
    k8s_values = Path(__file__).parent / "signoz-k8s-values.yaml"
    
    if not signoz_values.exists():
        print(f"❌ {signoz_values} 파일이 없습니다")
        sys.exit(1)
    print("✅ signoz-values.yaml 파일 존재")
    
    if not k8s_values.exists():
        print(f"❌ {k8s_values} 파일이 없습니다")
        sys.exit(1)
    print("✅ signoz-k8s-values.yaml 파일 존재")


def deploy_signoz():
    """SigNoz 애플리케이션 배포"""
    print("\n🚀 SigNoz 애플리케이션 배포 시작...")
    
    signoz_values = Path(__file__).parent / "signoz-values.yaml"
    
    try:
        # Helm 레포지토리 추가
        print("📦 SigNoz Helm 레포지토리 추가 중...")
        run_command(["helm", "repo", "add", "signoz", "https://charts.signoz.io"])
        
        # Helm 레포지토리 업데이트
        print("🔄 Helm 레포지토리 업데이트 중...")
        run_command(["helm", "repo", "update"])
        
        # SigNoz 애플리케이션 설치
        print(f"⚙️  SigNoz 애플리케이션을 {NAMESPACE} 네임스페이스에 설치 중...")
        print("⏳ 설치에는 몇 분이 소요될 수 있습니다...")
        
        install_cmd = [
            "helm", "upgrade", "--install", "signoz", "signoz/signoz",
            "--namespace", NAMESPACE,
            "--create-namespace",
            "--wait",
            "--timeout", "1h",
            "-f", str(signoz_values)
        ]
        
        run_command(install_cmd)
        
        print("✅ SigNoz 애플리케이션 설치 완료!")
        
    except subprocess.CalledProcessError as e:
        print(f"❌ SigNoz 애플리케이션 설치 실패: {e}")
        sys.exit(1)


def deploy_k8s_infra():
    """SigNoz K8s 인프라 (메트릭 수집) 배포"""
    print("\n📊 SigNoz K8s 인프라 배포 시작...")
    
    k8s_values = Path(__file__).parent / "signoz-k8s-values.yaml"
    
    try:
        # K8s 인프라 설치
        print(f"⚙️  K8s 메트릭 수집 에이전트를 {NAMESPACE} 네임스페이스에 설치 중...")
        print("⏳ DaemonSet 배포에는 1-2분이 소요될 수 있습니다...")
        
        install_cmd = [
            "helm", "upgrade", "--install", "k8s-infra", "signoz/k8s-infra",
            "--namespace", NAMESPACE,
            "--wait",
            "--timeout", "30m",
            "-f", str(k8s_values)
        ]
        
        run_command(install_cmd)
        
        print("✅ K8s 인프라 설치 완료!")
        
    except subprocess.CalledProcessError as e:
        print(f"❌ K8s 인프라 설치 실패: {e}")
        sys.exit(1)


def verify_deployment():
    """배포 상태 확인"""
    print("\n🔍 배포 상태 확인 중...")
    
    
    try:
        # Pod 상태 확인
        print(f"📊 {NAMESPACE} 네임스페이스의 Pod 상태:")
        run_command(["kubectl", "get", "pods", "-n", NAMESPACE, "-o", "wide"])
        
        # 서비스 상태 확인
        print(f"\n🌐 {NAMESPACE} 네임스페이스의 서비스 상태:")
        run_command(["kubectl", "get", "svc", "-n", NAMESPACE])
        
        # DaemonSet 상태 확인 (K8s 인프라)
        print(f"\n🔧 {NAMESPACE} 네임스페이스의 DaemonSet 상태:")
        run_command(["kubectl", "get", "daemonset", "-n", NAMESPACE])
        
        # Helm 릴리스 상태 확인
        print("\n📋 Helm 릴리스 상태:")
        run_command(["helm", "list", "-n", NAMESPACE])
        
        print("\n✅ 배포 상태 확인 완료!")
        
        # 접속 정보 안내
        print("\n📋 SigNoz 접속 정보:")
        print("1. 포트 포워딩을 통한 접속:")
        print(f"   kubectl port-forward -n {NAMESPACE} svc/signoz 28080:8080")
        print("2. 브라우저에서 http://localhost:28080 접속")
        print("3. 기본 로그인 정보:")
        print("   - Email: admin@signoz.io")
        print("   - Password: admin")
        
    except subprocess.CalledProcessError as e:
        print(f"⚠️  배포 상태 확인 중 오류: {e}")


def main():
    """메인 실행 함수"""
    print("🏠 홈서버 SigNoz 배포 스크립트")
    print("=" * 50)
    
    try:
        check_prerequisites()
        deploy_signoz()
        deploy_k8s_infra()
        verify_deployment()
        
        print("\n🎉 SigNoz 전체 스택 배포가 성공적으로 완료되었습니다!")
        print("💡 몇 분 후 SigNoz에서 호스트 및 K8s 메트릭을 확인할 수 있습니다.")
        
    except KeyboardInterrupt:
        print("\n❌ 사용자에 의해 배포가 중단되었습니다")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ 예상치 못한 오류 발생: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()