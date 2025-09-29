#!/usr/bin/env bash
set -euo pipefail

echo "🧹 Cleaning up leftover resources for SigNoz ..."

# 1. PVC + PV 정리
echo "👉 Cleaning PVC/PV..."
kubectl get pvc --all-namespaces | grep -i signoz | awk '{print $2" "$1}' | while read ns pvc; do
  echo "  - Deleting PVC $pvc in namespace $ns"
  kubectl patch pvc "$pvc" -n "$ns" -p '{"metadata":{"finalizers":[]}}' --type=merge || true
  kubectl delete pvc "$pvc" -n "$ns" --ignore-not-found
done

kubectl get pv | grep -i signoz | awk '{print $1}' | while read pv; do
  echo "  - Deleting PV $pv"
  kubectl patch pv "$pv" -p '{"metadata":{"finalizers":[]}}' --type=merge || true
  kubectl delete pv "$pv" --ignore-not-found
done

# 2. CRD 정리
echo "👉 Cleaning CRDs..."
kubectl get crd | grep -E 'signoz|clickhouse' | awk '{print $1}' | while read crd; do
  echo "  - Deleting CRD $crd"
  kubectl patch crd "$crd" -p '{"metadata":{"finalizers":[]}}' --type=merge || true
  kubectl delete crd "$crd" --ignore-not-found
done

# 3. ClusterRole / Binding 정리
echo "👉 Cleaning ClusterRoles/Bindings..."
kubectl get clusterrole | grep -E 'signoz|clickhouse' | awk '{print $1}' | while read cr; do
  echo "  - Deleting ClusterRole $cr"
  kubectl delete clusterrole "$cr" --ignore-not-found
done

kubectl get clusterrolebinding | grep -E 'signoz|clickhouse' | awk '{print $1}' | while read crb; do
  echo "  - Deleting ClusterRoleBinding $crb"
  kubectl delete clusterrolebinding "$crb" --ignore-not-found
done

# 4. Webhook 정리
echo "👉 Cleaning Webhooks..."
kubectl get validatingwebhookconfiguration | grep -E 'signoz|clickhouse' | awk '{print $1}' | while read vwh; do
  echo "  - Deleting ValidatingWebhook $vwh"
  kubectl delete validatingwebhookconfiguration "$vwh" --ignore-not-found
done

kubectl get mutatingwebhookconfiguration | grep -E 'signoz|clickhouse' | awk '{print $1}' | while read mwh; do
  echo "  - Deleting MutatingWebhook $mwh"
  kubectl delete mutatingwebhookconfiguration "$mwh" --ignore-not-found
done

# 5. Namespace 재확인
echo "👉 Cleaning leftover namespaces..."
kubectl get ns | grep -i signoz | awk '{print $1}' | while read ns; do
  echo "  - Deleting namespace $ns"
  kubectl patch ns "$ns" -p '{"spec":{"finalizers":[]}}' --type=merge || true
  kubectl delete ns "$ns" --ignore-not-found
done

echo "✅ SigNoz cleanup finished."