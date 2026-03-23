---
id: rbac-denied
category: rbac
severity: warning
title: "RBAC 权限拒绝"
related_states: [Forbidden, forbidden, denied, rbac, permission, unauthorized, accessdenied]
tags: [rbac, permission, denied, authorization]
---

## 故障现象
用户或 ServiceAccount 执行操作时被拒绝，收到 Forbidden 错误。
通常是 Role/ClusterRole 权限不足或 RoleBinding 未正确绑定。

## 典型症状
- 操作返回 403 Forbidden
- 错误信息包含 "disallowed by policy" 或 " forbidden"
- 特定 API 调用失败

## 排查步骤

1. 查看完整的错误信息：
   kubectl auth can-i <action> --as=<user-identifier>

2. 检查用户的 RoleBinding：
   kubectl get rolebinding -n <namespace>
   kubectl get clusterrolebinding

3. 检查 Role/ClusterRole 定义：
   kubectl get role -n <namespace>
   kubectl get clusterrole

4. 测试实际权限：
   kubectl auth can-i get pods --as=system:serviceaccount:<namespace>:<sa-name>

5. 查看认证用户：
   kubectl config current-context

6. 检查 ServiceAccount：
   kubectl get serviceaccount <sa-name> -n <namespace>

## 修复建议

1. 创建合适的 Role/ClusterRole
2. 创建 RoleBinding/ClusterRoleBinding 绑定权限
3. 检查操作是否需要特定的口令
4. 确认命名空间是否正确
5. 如果是 Pod 中的 ServiceAccount，检查是否正确挂载

## 可执行命令

```bash
# 测试权限
kubectl auth can-i get pods --as=<user-name>
kubectl auth can-i get pods --as=system:serviceaccount:<namespace>:<sa-name>

# 查看 RoleBinding
kubectl get rolebinding -n <namespace> -o yaml

# 查看 ClusterRoleBinding
kubectl get clusterrolebinding -o yaml

# 列出用户的权限
kubectl auth can-i --list --as=<user-name>

# 检查 ServiceAccount
kubectl get serviceaccount <sa-name> -n <namespace> -o yaml

# 创建 Role 示例
kubectl create role <role-name> --verb=get,list --resource=pods -n <namespace>

# 创建 RoleBinding
kubectl create rolebinding <name> --role=<role-name> --user=<user-name> -n <namespace>
```
