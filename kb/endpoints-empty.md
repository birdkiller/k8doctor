---
id: endpoints-empty
category: network
severity: warning
title: "Endpoints 为空"
related_states: [Endpoints, endpoints, empty, noendpoints, nobackend]
tags: [service, endpoint, selector, backend]
---

## 故障现象
Service 的 Endpoints 为空，没有后端 Pod 接收流量。
访问 Service 时会直接被拒绝或无响应。

## 典型症状
- kubectl get endpoints 显示 ENDPOINTS 为空
- Service 存在但无法转发流量
- 应用日志显示无请求到达

## 排查步骤

1. 确认 Service 存在：
   kubectl get svc <service-name> -n <namespace>

2. 查看 Endpoints：
   kubectl get endpoints <service-name> -n <namespace>

3. 检查 Pod 是否 Running：
   kubectl get pods -n <namespace>

4. 检查 Selector 匹配：
   serviceSelector=$(kubectl get svc <service-name> -n <namespace> -o jsonpath='{.spec.selector}')
   kubectl get pods -n <namespace> -l "$serviceSelector"

5. 检查 Pod 端口配置：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.containers[*].ports}'

6. 检查 Pod 是否通过 ReadinessProbe：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 5 "Readiness"

## 修复建议

1. 确保后端 Pod 处于 Running 且 Ready
2. 修正 Service Selector 确保匹配正确的 Pod
3. 检查 targetPort 配置是否正确
4. 如果 Pod 太多，排查是否所有 Pod 都被 ReadinessProbe 拒绝

## 可执行命令

```bash
# 查看 Endpoints
kubectl get endpoints <service-name> -n <namespace>

# 获取 Service Selector
kubectl get svc <service-name> -n <namespace> -o jsonpath='{.spec.selector}' | jq

# 查看匹配到的 Pod
kubectl get pods -n <namespace> -l "$(kubectl get svc <service-name> -n <namespace> -o jsonpath='{.spec.selector}')"

# 查看 Pod 状态
kubectl get pods -n <namespace> -o wide

# 检查 Pod ReadinessProbe
kubectl describe pod <pod-name> -n <namespace> | grep -A 10 "Readiness"

# 如果使用 Deployment，检查 ReplicaSet
kubectl get rs -n <namespace> -l "$(kubectl get svc <service-name> -n <namespace> -o jsonpath='{.spec.selector}')"
```
