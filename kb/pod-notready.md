---
id: pod-notready
category: pod
severity: warning
title: "Pod Running 但 NotReady"
related_states: [NotReady, notready, running, notReady]
tags: [pod, readiness, probe, network]
---

## 故障现象
Pod 处于 Running 状态，但 Ready 不为 1/1，表示健康检查未通过。
服务无法接收流量，但进程仍在运行。

## 典型症状
- Pod 状态：Running
- Ready 显示：0/1 或类似
- Readiness Probe 失败
- 应用可能还在运行但无法响应请求

## 排查步骤

1. 查看 Pod 详细信息：
   kubectl describe pod <pod-name> -n <namespace>

2. 检查 Readiness Probe 配置：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.containers[*].readinessProbe}'

3. 查看探测日志：
   kubectl logs <pod-name> -n <namespace>

4. 手动测试健康检查端点：
   kubectl exec <pod-name> -n <namespace> -- curl -k localhost:<port>/health

5. 检查网络策略是否阻止了探测：
   kubectl get networkpolicy -n <namespace>

6. 检查服务 Endpoints：
   kubectl get endpoints <service-name> -n <namespace>

## 修复建议

1. 检查应用健康检查端点是否正常
2. 调整 readinessProbe 延迟和超时时间
3. 确认容器端口配置正确
4. 检查网络策略

## 可执行命令

```bash
# 查看 Ready 状态
kubectl get pod <pod-name> -n <namespace> -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,READY:.status.conditions[?(@.type=='Ready')].status

# 查看探测详情
kubectl describe pod <pod-name> -n <namespace> | grep -A 10 "Readiness"

# 手动测试健康端点
kubectl exec <pod-name> -n <namespace> -- wget -qO- localhost:<port>/health

# 查看 Endpoints
kubectl get endpoints <service-name> -n <namespace> -o yaml

# 检查服务选择器
kubectl get svc <service-name> -n <namespace> -o jsonpath='{.spec.selector}'
```
