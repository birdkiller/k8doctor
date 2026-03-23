---
id: node-diskpressure
category: node
severity: warning
title: "Node DiskPressure"
related_states: [DiskPressure, diskpressure, disk, diskfull, diskquota]
tags: [node, disk, pressure, storage]
---

## 故障现象
Node 磁盘空间不足或 inode 耗尽，触发 Kubelet 磁盘压力驱逐。
容器无法创建，日志无法写入，新的 Pod 无法调度。

## 典型症状
- Node 状态包含 DiskPressure=True
- Pod 创建失败或被驱逐
- kubelet 日志出现 "no space left on device"
- docker/containerd 存储持续增长

## 排查步骤

1. 检查 Node 磁盘状态：
   kubectl describe node <node-name> | grep -A 10 "Conditions"

2. 检查磁盘使用：
   df -h
   df -i

3. 查看 Docker 存储使用：
   docker system df

4. 检查大文件：
   du -sh /var/log/*
   find /var/log -name "*.log" -size +100M

5. 检查容器日志：
   docker logs <container-id>

## 修复建议

1. 清理 Docker 存储：
   docker system prune -a -f --volumes

2. 清理日志文件：
   journalctl --vacuum-size=100M
   find /var/log -name "*.gz" -mtime +7 -delete

3. 调整 kubelet 驱逐阈值
4. 配置日志轮转
5. 考虑增加节点磁盘或挂载新磁盘

## 可执行命令

```bash
# 查看磁盘使用
df -h / /var/lib/docker

# 检查 inode
df -i

# 查看 Docker 存储占用
docker system df -v

# 清理 Docker
docker system prune -a -f --volumes

# 清理 journal 日志
journalctl --vacuum-size=100M

# 查看大日志文件
find /var/log -name "*.log" -exec ls -lh {} \; | sort -k5 -hr | head -10

# 限制 Docker 存储
cat /etc/docker/daemon.json
# 添加 { "storage-driver": "overlay2", "data-root": "/path/to/disk" }
```
