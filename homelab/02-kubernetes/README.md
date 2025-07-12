# Kubernetes 앱 배포

Kubernetes 클러스터 애플리케이션 배포 및 관리

## 배포 가능한 애플리케이션

각 애플리케이션별 세부 설정은 해당 디렉토리의 배포 스크립트에서 확인하세요.

### SigNoz 모니터링 스택

**기술 스택**: Helm, OpenTelemetry, ClickHouse  
**기능**: 호스트 메트릭, K8s 메트릭 수집 및 모니터링

```bash
cd signoz
python3 deploy.py
```
