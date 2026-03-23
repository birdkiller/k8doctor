---
id: ingress-502
category: network
severity: warning
title: "Ingress 502/504 错误"
related_states: [502, 504, BadGateway, GatewayTimeout, ingresserror, 502error]
tags: [ingress, nginx, proxy, gateway, 502]
---

## 故障现象
通过 Ingress 访问服务返回 502 Bad Gateway 或 504 Gateway Timeout。
说明 Ingress Controller 无法连接到后端服务。

## 典型症状
- HTTP 状态码 502 或 504
- Nginx Ingress 返回 502
- 错误信息：no live upstream / upstream prematurely closed
- 直接访问 Service 正常

## 排查步骤

1. 检查 Ingress 是否配置正确：
   kubectl get ingress <ingress-name> -n <namespace>

2. 检查 Ingress Controller 是否运行：
   kubectl get pods -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx

3. 查看 Ingress 事件：
   kubectl describe ingress <ingress-name> -n <namespace>

4. 检查后端 Service 和 Endpoints：
   kubectl get svc -n <namespace>
   kubectl get endpoints <backend-service> -n <namespace>

5. 检查 Ingress Controller 日志：
   kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx --tail=100

6. 验证 Ingress Class 配置：
   kubectl get ingressclass

## 修复建议

1. 确保后端 Service 存在且有 Endpoints
2. 检查 Service Port 和 Ingress Port 配置是否匹配
3. 确保 Ingress 绑定了正确的 IngressClass
4. 检查 TLS 证书是否有效
5. 增加 Ingress Controller 副本

## 可执行命令

```bash
# 查看 Ingress 详情
kubectl describe ingress <ingress-name> -n <namespace>

# 查看 Ingress Controller 状态
kubectl get pods -n ingress-nginx

# 查看 Ingress Controller 日志
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx --tail=200

# 测试后端 Service
kubectl exec -it <test-pod> -n <namespace> -- curl -v <service-name>:<port>

# 查看 Endpoints
kubectl get endpoints <backend-service> -n <namespace>

# 检查 IngressClass
kubectl get ingressclass
kubectl get ingress <ingress-name> -n <namespace> -o jsonpath='{.spec.ingressClassName}'

# 查看配置映射
kubectl get configmap -n ingress-nginx
```
