---
id: node-notready
category: node
severity: critical
title: "Node NotReady"
related_states: [NotReady, notready, nodered, unreachable, nodemachinedown]
tags: [node, notready, cluster, kubelet]
---

## 故障现象
Node 状态变为 NotReady，不再接受新的 Pod 调度。
运行在该节点上的 Pod 可能受影响，集群容量下降。

## 典型症状
- Node 状态：NotReady
- Ready 条件为 False 或 Unknown
- kubelet 可能无法正常通信
- 运行中的 Pod 状态不变但无法管理

## 排查步骤

1. 检查 Node 详细状态：
   kubectl describe node <node-name>

2. 检查 kubelet 日志：
   journalctl -u kubelet -n 100 --no-pager | grep <node-name>

3. 验证 kubelet 服务状态：
   systemctl status kubelet

4. 检查 Docker/containerd 服务：
   systemctl status docker
   systemctl status containerd

5. 检查网络连通性：
   ping <node-ip>
   telnet <node-ip> 10250

6. 检查磁盘空间：
   df -h

7. 检查内存使用：
   free -m

8. 检查 CPU 使用：
   top -bn1 | head -20

## 修复建议

1. 如果是临时性问题，等待节点自动恢复
2. 如果 kubelet 卡死，重启 kubelet 服务
3. 清理磁盘空间释放 /var/lib/docker
4. 如果节点无法恢复，可将其从集群中移除：
   kubectl drain <node-name> --ignore-daemonsets --delete-local-data
   kubectl delete node <node-name>

## 可执行命令

```bash
# 查看 Node 状态
kubectl get node <node-name> -o wide

# 查看 NotReady 原因
kubectl describe node <node-name> | grep -A 5 "Conditions"

# 查看 kubelet 日志
journalctl -u kubelet -n 200 --no-pager

# 检查 kubelet 是否运行
systemctl status kubelet

# 重启 kubelet
systemctl restart kubelet

# 检查 Docker
systemctl status docker

# 驱逐节点（谨慎操作）
kubectl drain <node-name> --ignore-daemonsets --delete-local-data --force

# 从集群移除节点
kubectl delete node <node-name>
```
