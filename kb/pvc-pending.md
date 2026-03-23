---
id: pvc-pending
category: storage
severity: warning
title: "PVC Pending"
related_states: [Pending, pvc, pvcpending, claimpending, waitingforclaim]
tags: [pvc, storage, pending, claim, volume]
---

## 故障现象
PersistentVolumeClaim 一直处于 Pending 状态，Pod 无法挂载存储。
通常是因为没有合适的 PV 或 StorageClass 配置问题。

## 典型症状
- kubectl get pvc 显示 STATUS=Pending
- Pod 无法启动，报 MountVolume error
- PVC 绑定失败

## 排查步骤

1. 查看 PVC 详情：
   kubectl describe pvc <pvc-name> -n <namespace>

2. 检查 StorageClass 是否存在：
   kubectl get storageclass

3. 查看 PV 状态：
   kubectl get pv

4. 检查 PVC 请求的 StorageClass：
   kubectl get pvc <pvc-name> -n <namespace> -o jsonpath='{.spec.storageClassName}'

5. 如果使用动态供给，检查 Provisioner：
   kubectl get pods -n kube-system | grep provisioner

6. 查看相关的事件：
   kubectl get events -n <namespace> | grep -i pvc

## 修复建议

1. 确保 StorageClass 存在且配置正确
2. 如果是手动供给，创建匹配的 PV
3. 检查集群是否有足够的存储
4. 检查 Provisioner 是否正常运行
5. 确认 PVC 的 accessModes 与 PV 匹配

## 可执行命令

```bash
# 查看 PVC 详情
kubectl describe pvc <pvc-name> -n <namespace>

# 查看 PVC 状态
kubectl get pvc -n <namespace>

# 查看 StorageClass
kubectl get storageclass

# 查看 PV
kubectl get pv

# 查看动态Provisioner
kubectl get pods -n kube-system | grep -i provisioner

# 查看相关事件
kubectl get events -n <namespace> --sort-by='.lastTimestamp' | grep -i pvc

# 如果需要手动创建 PV：
# kubectl create -f pv.yaml
```
