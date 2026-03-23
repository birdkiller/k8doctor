---
id: pod-imagepullbackoff
category: pod
severity: warning
title: "Pod ImagePullBackOff"
related_states: [ImagePullBackOff, imagepull, pullfailed, 401, 403, 404]
tags: [pod, image, pull, registry]
---

## 故障现象
Pod 无法拉取容器镜像，处于 ImagePullBackOff 状态。
通常是因为镜像不存在、权限不足或网络问题。

## 典型症状
- Pod 状态：ImagePullBackOff
- 错误信息包含 "ImagePullBackOff" 或 "rpc error"
- Backoff 倒计时持续增长

## 排查步骤

1. 查看 Pod 详情获取具体错误：
   kubectl describe pod <pod-name> -n <namespace>

2. 检查镜像名称是否正确：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.containers[*].image}'

3. 验证镜像是否存在于仓库：
   kubectl run test --image=<image-name> --rm -it -- /bin/sh

4. 检查 Secret 是否正确配置（私有仓库）：
   kubectl get secret <secret-name> -n <namespace>

5. 检查 ServiceAccount 是否绑定了正确的 imagePullSecret：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.imagePullSecrets}'

## 修复建议

1. 确认镜像名称和标签正确
2. 创建或更新 imagePullSecrets：
   kubectl create secret docker-registry <secret-name> \
     --docker-server=<registry> \
     --docker-username=<user> \
     --docker-password=<pass> \
     -n <namespace>
3. 将 secret 绑定到 ServiceAccount
4. 检查网络策略是否允许访问外网/镜像仓库

## 可执行命令

```bash
# 查看拉取失败详情
kubectl describe pod <pod-name> -n <namespace> | grep -A 5 "Events"

# 检查镜像
crictl images | grep <image-name>

# 手动测试拉取镜像
crictl pull <image-name>

# 创建私有仓库 Secret
kubectl create secret docker-registry <secret-name> \
  --docker-server=https://index.docker.io/v1/ \
  --docker-username=<username> \
  --docker-password=<password> \
  --docker-email=<email> \
  -n <namespace>

# 绑定到 ServiceAccount
kubectl patch serviceaccount default -n <namespace> \
  -p '{"imagePullSecrets":[{"name":"<secret-name>"}]}'
```
