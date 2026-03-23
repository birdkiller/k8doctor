---
id: pod-pending
category: pod
severity: warning
title: "Pod Pending"
related_states: [Pending, scheduling, waiting, nomogeneous]
tags: [pod, pending, schedule, resources]
---

## 故障现象
Pod 长时间处于 Pending 状态，无法被调度到任何 Node 上。
通常是资源不足、亲和性限制或调度器问题导致。

## 典型症状
- Pod 状态：Pending
- 没有任何 Node 被分配（Node Selector 为空）
- Events 中显示 "FailedScheduling"

## 排查步骤

1. 查看 Pod 调度详情：
   kubectl describe pod <pod-name> -n <namespace>

2. 检查资源请求是否超过集群总量：
   kubectl describe node | grep -A 5 "Allocated resources"

3. 查看 Node 资源使用情况：
   kubectl top nodes

4. 检查污点和容忍：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 10 "Tolerations"

5. 检查亲和性规则：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 20 "Affinity"

6. 验证 PVC 是否已绑定：
   kubectl get pvc -n <namespace>

## 修复建议

1. 增加集群节点或清理无用 Pod 释放资源
2. 调整 Pod 的 resources.requests 值
3. 添加合适的容忍来匹配 Node 污点
4. 修正亲和性规则
5. 确保 PVC 正常绑定

## 可执行命令

```bash
# 查看调度失败原因
kubectl describe pod <pod-name> -n <namespace> | grep -A 10 "Events"

# 查看所有 Node 资源
kubectl describe node | grep -A 8 "Allocated resources"

# 查看 CPU/Memory 压力节点
kubectl top nodes

# 查看 Pod 资源请求
kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.containers[*].resources.requests}'

# 检查污点
kubectl get node -o custom-columns=NAME:.metadata.name,TAINTS:.spec.taints

# 临时禁用调度保护测试
kubectl patch node <node-name> -p '{"spec":{"unschedulable":false}}'
```
