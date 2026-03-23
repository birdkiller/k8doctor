---
id: pod-crashloop
category: pod
severity: critical
title: "Pod CrashLoopBackOff"
related_states: [CrashLoopBackOff, crashloop, restart, restarting]
tags: [pod, crash, restart, loop]
---

## 故障现象
Pod 处于 CrashLoopBackOff 状态，持续重启失败。
容器启动后立即退出，然后重新启动，形成恶性循环。

## 典型症状
- Pod 状态：CrashLoopBackOff
- Restart Count 持续增加
- Ready 状态不稳定
- 日志中可以看到应用报错信息

## 排查步骤

1. 查看 Pod 详细状态：
   kubectl describe pod <pod-name> -n <namespace>

2. 查看容器日志：
   kubectl logs <pod-name> -n <namespace> --previous

3. 检查容器启动命令和参数：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.containers[*].command}'

4. 检查镜像是否存在：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.containers[*].image}'

5. 检查 ConfigMap/Secret 挂载是否正确：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 10 "Mounts"

## 修复建议

1. 根据错误日志修复应用配置
2. 检查环境变量是否正确设置
3. 确认依赖服务是否可达
4. 检查数据目录权限

## 可执行命令

```bash
# 查看 CrashLoopBackOff 详情
kubectl describe pod <pod-name> -n <namespace> | grep -E "(State:|Last State:|Exit Code:)"

# 查看最近日志
kubectl logs <pod-name> -n <namespace> --tail=100

# 查看上次失败的容器日志
kubectl logs <pod-name> -n <namespace> --previous

# 检查镜像是否存在
kubectl run test --image=<image-name> --rm -it -- /bin/sh
```
