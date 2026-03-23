---
id: service-unreachable
category: network
severity: warning
title: "Service 无法访问"
related_states: [unreachable, connectionrefused, timeout, refused, refusedconnection]
tags: [service, network, unreachable, connection]
---

## 故障现象
无法通过 Service 访问应用，连接超时或被拒绝。
但直接通过 Pod IP 可以访问，说明网络层面有问题。

## 典型症状
- curl/wget 访问 Service IP 超时
- 连接被拒绝 (Connection Refused)
- DNS 解析正常但无法连接
- 直接访问 Pod 正常

## 排查步骤

1. 检查 Service 是否存在：
   kubectl get svc <service-name> -n <namespace>

2. 查看 Service 的 Endpoints：
   kubectl get endpoints <service-name> -n <namespace>

3. 检查 Selector 是否匹配：
   kubectl get svc <service-name> -n <namespace> -o jsonpath='{.spec.selector}'
   kubectl get pods -n <namespace> -l <selector>

4. 检查 Pod 是否 Running：
   kubectl get pods -n <namespace> -l <selector>

5. 测试 DNS 解析：
   kubectl exec <pod-name> -n <namespace> -- nslookup <service-name>

6. 检查网络策略：
   kubectl get networkpolicy -n <namespace>

## 修复建议

1. 确保 Pod Selector 正确匹配
2. 确保 Pod 处于 Running 状态
3. 检查 Service 端口配置
4. 调整网络策略允许流量
5. 检查 kube-proxy 是否正常运行

## 可执行命令

```bash
# 查看 Service 和 Endpoints
kubectl get svc <service-name> -n <namespace>
kubectl get endpoints <service-name> -n <namespace>

# 测试 Service 连接
kubectl exec <pod-name> -n <namespace> -- curl -v <service-ip>:<port>

# 测试 DNS
kubectl exec <pod-name> -n <namespace> -- nslookup <service-name>
kubectl exec <pod-name> -n <namespace> -- cat /etc/resolv.conf

# 检查 kube-proxy
kubectl get pods -n kube-system -l k8s-app=kube-proxy
kubectl logs -n kube-system -l k8s-app=kube-proxy --tail=100

# 查看网络策略
kubectl get networkpolicy -n <namespace>

# 检查 iptables 规则
iptables -L -n -t nat | grep <service-name>
```
