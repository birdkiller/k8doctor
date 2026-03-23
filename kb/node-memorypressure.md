---
id: node-memorypressure
category: node
severity: warning
title: "Node MemoryPressure"
related_states: [MemoryPressure, memorypressure, memory, mempressure]
tags: [node, memory, pressure, resource]
---

## 故障现象
Node 报告内存压力，触发 Kubelet 内存压力驱逐。
新的 Pod 可能无法调度到该节点，已有的 Pod 可能被驱逐。

## 典型症状
- Node 状态包含 MemoryPressure=True
- Events 中出现 "MemoryPressure" 相关事件
- Pod 被频繁驱逐
- kubelet 开始驱逐 Pod 以释放内存

## 排查步骤

1. 检查 Node 内存状态：
   kubectl describe node <node-name> | grep -A 10 "MemoryPressure"

2. 查看 Node 内存分配：
   kubectl describe node <node-name> | grep -A 5 "Allocated resources"

3. 检查系统内存使用：
   free -m
   cat /proc/meminfo

4. 查看占用内存最多的进程：
   ps aux --sort=-%mem | head -10

5. 查看被驱逐的 Pod：
   kubectl get events -n <namespace> | grep -i evicted

## 修复建议

1. 增加集群节点或迁移部分 Pod
2. 清理不必要的进程或容器
3. 调整 kubelet --eviction-hard 参数
4. 设置合适的资源请求和限制让调度更合理
5. 考虑增加 Node 的物理内存

## 可执行命令

```bash
# 查看 Node 内存条件
kubectl describe node <node-name> | grep -A 5 "Conditions"

# 查看 Node 内存分配详情
kubectl describe node <node-name> | grep -A 8 "Allocated resources"

# 检查系统内存
free -h

# 查看内存使用最多的容器
docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}" | sort -k2 -hr | head -10

# crictl 查看
crictl stats | sort -k4 -hr | head -10

# 清理已停止容器
docker container prune -f

# 清理无用镜像
docker image prune -a -f
```
