---
id: dns-failure
category: network
severity: warning
title: "DNS 解析失败"
related_states: [dns, DNS, nxdomain, notfound, resolv, nameserver]
tags: [dns, coredns, kube-dns, resolution]
---

## 故障现象
Pod 内无法解析集群内部或外部域名。
curl 使用 Service 名称访问失败，nslookup 超时。

## 典型症状
- nslookup 失败
- curl 使用 Service 名称无法连接
- DNS 查询超时或返回 SERVFAIL
- /etc/resolv.conf 配置正确但无法使用

## 排查步骤

1. 测试 DNS 解析：
   kubectl exec <pod-name> -n <namespace> -- nslookup kubernetes.default
   kubectl exec <pod-name> -n <namespace> -- cat /etc/resolv.conf

2. 检查 CoreDNS/CoreDNS 是否运行：
   kubectl get pods -n kube-system -l k8s-app=kube-dns

3. 查看 CoreDNS 日志：
   kubectl logs -n kube-system -l k8s-app=kube-dns --tail=100

4. 测试 DNS 缓存：
   kubectl exec <pod-name> -n <namespace> -- nslookup www.baidu.com

5. 检查网络连通性：
   kubectl exec <pod-name> -n <namespace> -- ping 8.8.8.8

6. 检查 Node 的 resolv.conf：
   kubectl exec <pod-name> -n <namespace> -- cat /etc/resolv.conf

## 修复建议

1. 重启 CoreDNS：
   kubectl rollout restart deployment coredns -n kube-system

2. 检查 nodelocaldns 配置
3. 确认网络插件(CNI)正常工作
4. 检查是否有网络策略阻止 DNS 流量
5. 验证集群 DNS 配置：
   kubectl config current-context
   kubectl cluster-info

## 可执行命令

```bash
# 测试集群 DNS
kubectl exec -it <pod-name> -n <namespace> -- nslookup kubernetes.default

# 测试外部 DNS
kubectl exec -it <pod-name> -n <namespace> -- nslookup www.baidu.com

# 查看 CoreDNS 状态
kubectl get pods -n kube-system -l k8s-app=kube-dns

# 查看 CoreDNS 日志
kubectl logs -n kube-system -l k8s-app=kube-dns --tail=200

# 重启 CoreDNS
kubectl rollout restart deployment coredns -n kube-system

# 检查 DNS 配置
kubectl exec -it <pod-name> -n <namespace> -- cat /etc/resolv.conf

# 检查网络策略
kubectl get networkpolicy --all-namespaces
```
