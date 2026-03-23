---
id: pod-oomkilled
category: pod
severity: critical
title: "Pod 发生 OOMKilled"
related_states: [OOMKilled, OOM, killed, 137, exit137]
tags: [memory, pod, crash, oom, limit]
---

## 故障现象
Pod 被强制终止，退出码为 137，表示收到 SIGKILL 信号。
通常原因是容器内存使用量超过 limits 限制。

## 典型症状
- Pod 状态：OOMKilled
- 重启次数持续增加
- Exit Code: 137
- 日志中无错误输出（进程被直接 kill）
- kubectl describe 显示 Last State 为 Terminated

## 排查步骤

1. 查看 Pod 详情，重点关注 Last State：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 10 "Last State"

2. 查看内存使用量（需 metrics-server）：
   kubectl top pod <pod-name> -n <namespace>

3. 检查 limits 配置：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.containers[*].resources}'

4. 查看历史资源使用：
   kubectl top pod <pod-name> -n <namespace> --历史

5. 检查 Node 是否有内存压力：
   kubectl describe node <node-name> | grep -A 5 "MemoryPressure"

## 修复建议

1. 调整 memory limits：建议设为当前稳定使用量的 1.5-2 倍
2. 检查 Java/Node.js 等运行时 heap 设置是否合理
3. 使用 profiling 工具定位内存泄漏
4. 考虑增加 Pod 副本分散负载

## 可执行命令

```bash
# 查看 Pod 状态和重启原因
kubectl describe pod <pod-name> -n <namespace> | grep -A 10 "Last State"

# 查看内存使用
kubectl top pod <pod-name> -n <namespace>

# 查看当前 limits 配置
kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.containers[*].resources.limits}'

# 快速扩容临时恢复
kubectl scale deployment <deployment-name> --replicas=0 -n <namespace> && \
kubectl scale deployment <deployment-name> --replicas=3 -n <namespace>
```
