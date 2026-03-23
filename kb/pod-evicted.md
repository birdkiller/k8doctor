---
id: pod-evicted
category: pod
severity: warning
title: "Pod 被 Evicted"
related_states: [Evicted, evicted, preempted, evictedby]
tags: [pod, evicted, resources, pressure]
---

## 故障现象
Pod 被 Kubernetes 驱逐，通常是因为资源压力、节点维护或驱逐策略触发。
Pod 状态显示为 Evicted，容器被强制终止。

## 典型症状
- Pod 状态：Evicted
- 退出码可能为 137 或其他
- 通常发生在资源紧张或节点维护期间
- Deployment 会尝试重新创建 Pod

## 排查步骤

1. 查看 Pod 驱逐详情：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 15 "Events"

2. 检查 Node 资源状态：
   kubectl describe node <node-name> | grep -A 10 "Conditions"

3. 查看是否有污点：
   kubectl describe node <node-name> | grep Taints

4. 检查集群整体资源：
   kubectl top nodes
   kubectl top pods -n <namespace>

5. 查看是否有 PodDisruptionBudget 干扰：
   kubectl get pdb -n <namespace>

## 修复建议

1. 增加集群容量或清理低优先级 Pod
2. 检查是否有节点维护操作在进行
3. 调整 Pod 资源请求和限制
4. 设置合理的 PriorityClass
5. 检查驱逐策略配置

## 可执行命令

```bash
# 查看驱逐原因
kubectl describe pod <pod-name> -n <namespace> | grep -A 5 "Events"

# 查看 Node 资源压力
kubectl describe node <node-name> | grep -A 10 "Conditions"

# 查看资源使用
kubectl top node <node-name>

# 查看所有 evicted Pod
kubectl get pods --all-namespaces | grep Evicted

# 清理 Evicted Pod
kubectl delete pods -n <namespace> --field-selector=status.phase=Evicted

# 检查 Pod 优先级
kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.priorityClassName}'
```
