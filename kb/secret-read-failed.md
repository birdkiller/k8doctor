---
id: secret-read-failed
category: rbac
severity: warning
title: "Secret 读取失败"
related_states: [secret, secrets, forbidden, denied, unable, readsecret]
tags: [secret, permission, config, secretaccess]
---

## 故障现象
Pod 无法读取 Secret，或者创建 Secret 时报错。
通常是 RBAC 权限问题或 Secret 不存在。

## 典型症状
- Pod 启动失败，报无法挂载 Secret
- kubectl get secret 报 Forbidden
- Secret 挂载的值为空
- 错误信息：Unable to mount secrets

## 排查步骤

1. 检查 Secret 是否存在：
   kubectl get secret <secret-name> -n <namespace>

2. 检查 Pod 使用的 ServiceAccount：
   kubectl get pod <pod-name> -n <namespace> -o jsonpath='{.spec.serviceAccountName}'

3. 测试 Secret 读取权限：
   kubectl auth can-i get secrets --as=system:serviceaccount:<namespace>:<sa-name>

4. 查看 Pod 事件：
   kubectl describe pod <pod-name> -n <namespace> | grep -A 10 Events

5. 检查 Secret 类型和数据：
   kubectl get secret <secret-name> -n <namespace> -o yaml

6. 检查 Role/RoleBinding：
   kubectl get role -n <namespace>
   kubectl get rolebinding -n <namespace>

## 修复建议

1. 创建 Secret：
   kubectl create secret generic <secret-name> \
     --from-literal=key=value \
     -n <namespace>

2. 创建 RoleBinding 授予读取权限：
   kubectl create rolebinding <name> \
     --role=secret-reader \
     --serviceaccount=<namespace>:<sa-name> \
     -n <namespace>

3. 或者使用默认的 edit 角色（测试环境）：
   kubectl edit rolebinding -n <namespace>

## 可执行命令

```bash
# 查看 Secret
kubectl get secret <secret-name> -n <namespace>
kubectl describe secret <secret-name> -n <namespace>

# 测试权限
kubectl auth can-i get secrets --as=system:serviceaccount:<namespace>:<sa-name>
kubectl auth can-i list secrets --as=system:serviceaccount:<namespace>:<sa-name>

# 创建 Secret
kubectl create secret generic <secret-name> \
  --from-literal=username=admin \
  --from-literal=password=secret \
  -n <namespace>

# 从文件创建
kubectl create secret generic <secret-name> \
  --from-file=./config.json \
  -n <namespace>

# 创建只读 Role
kubectl create role secret-reader \
  --verb=get,list,watch \
  --resource=secrets \
  -n <namespace>

# 绑定到 ServiceAccount
kubectl create rolebinding <name> \
  --role=secret-reader \
  --serviceaccount=<namespace>:<sa-name> \
  -n <namespace>
```
