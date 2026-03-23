---
id: storage-mount-failed
category: storage
severity: warning
title: "存储挂载失败"
related_states: [mount, mounting, notmount, volume, unmount, mountfailed, mounterror]
tags: [volume, mount, storage, pod]
---

## 故障现象
Pod 启动时无法挂载 Volume，容器无法启动。
可能是因为 PVC 未绑定、存储驱动问题或权限问题。

## 典型症状
- Pod 状态：ContainerCreating 或 Pending
- 事件：Failed to mount volume
- PVC 状态异常
- 挂载路径不存在或权限错误

## 排查步骤

1. 查看 Pod 挂载详情：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 20 "Volumes"

2. 检查 PVC 状态：
   kubectl get pvc -n <namespace>

3. 查看 Pod 事件：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 10 Events

4. 检查 Volume 插件是否支持：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.volumes[*].flexVolume}'

5. 检查 NFS/CIFS 等外部存储连通性：
   mount -t nfs <nfs-server>:/export /mnt

6. 检查 Node 上的存储驱动：
   ls /var/lib/kubelet/volumeplugins/

## 修复建议

1. 确保 PVC 已成功绑定
2. 检查存储类(StorageClass)配置是否正确
3. 确认 NFS/Ceph 等外部存储服务正常
4. 检查 Pod 的 fsGroup 配置
5. 验证 Volume 权限(AccessModes)

## 可执行命令

```bash
# 查看 Pod 挂载详情
kubectl describe pod <pod-name> -n <namespace> | grep -A 15 "Volumes"

# 查看 PVC
kubectl get pvc -n <namespace>

# 查看 PVC 详情
kubectl describe pvc <pvc-name> -n <namespace>

# 查看 Pod 事件
kubectl describe pod <pod-name> -n <namespace> | grep -A 10 Events

# 查看 PV
kubectl get pv
kubectl describe pv <pv-name>

# 检查 Node 上的存储插件
ls -la /var/lib/kubelet/volumeplugins/

# 测试 NFS 挂载
mount -t nfs <nfs-server>:/ /mnt

# 检查权限
kubectl exec -it <pod-name> -n <namespace> -- ls -la <mount-path>

# 检查 fsGroup
kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.securityContext.fsGroup}'
```
