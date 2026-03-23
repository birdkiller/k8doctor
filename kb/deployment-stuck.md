---
id: deployment-stuck
category: deployment
severity: warning
title: "Deployment 卡在 Progressing"
related_states: [Progressing, stuck, progressing, norolllout, rolloutstuck]
tags: [deployment, rollout, stuck, progress]
---

## 故障现象
Deployment 升级或部署卡住，一直处于 Progressing 状态。
新 ReplicaSet 无法达到期望的副本数，或者卡在某个版本。

## 典型症状
- Deployment status 显示 Progressing=True
- Available Replicas 小于 Desired
- 新 ReplicaSet 一直是 0 ready
-  rollout status 卡住不动

## 排查步骤

1. 查看 Deployment 状态：
   kubectl describe deployment <name> -n <namespace>

2. 检查 ReplicaSet 状态：
   kubectl get rs -n <namespace>

3. 查看新 Pod 的状态：
   kubectl get pods -n <namespace> -l <new-label-selector>

4. 检查滚动更新策略：
   kubectl get deployment <name> -n <namespace> -o jsonpath='{.spec.strategy}'

5. 查看 rollout 历史：
   kubectl rollout history deployment <name> -n <namespace>

6. 检查 Deployment 事件：
   kubectl describe deployment <name> -n <namespace> | grep -A 10 Events

## 修复建议

1. 如果 Pod 无法创建，参考 Pod 故障排查
2. 检查滚动更新参数：
   - spec.strategy.rollingUpdate.maxUnavailable
   - spec.strategy.rollingUpdate.maxSurge
3. 暂停滚动更新：
   kubectl rollout pause deployment <name> -n <namespace>
4. 恢复滚动更新：
   kubectl rollout resume deployment <name> -n <namespace>
5. 回滚到上一版本：
   kubectl rollout undo deployment <name> -n <namespace>

## 可执行命令

```bash
# 查看 Deployment 状态
kubectl describe deployment <name> -n <namespace>

# 查看 rollout 状态
kubectl rollout status deployment <name> -n <namespace>

# 查看 ReplicaSet
kubectl get rs -n <namespace>

# 查看新旧 Pod
kubectl get pods -n <namespace> -l deployment=<name>

# 暂停滚动更新
kubectl rollout pause deployment <name> -n <namespace>

# 恢复滚动更新
kubectl rollout resume deployment <name> -n <namespace>

# 回滚到上一版本
kubectl rollout undo deployment <name> -n <namespace>

# 回滚到指定版本
kubectl rollout undo deployment <name> --to-revision=<n> -n <namespace>

# 查看 rollout 历史
kubectl rollout history deployment <name> -n <namespace>
```
