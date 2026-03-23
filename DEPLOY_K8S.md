# K8s Doctor 容器化部署指南

将 K8s Doctor 部署到目标集群，实现故障自诊断。

## 快速开始

### 1. 构建 Docker 镜像

```bash
cd D:\workspace\k8doctor

# 构建镜像
docker build -t k8doctor:latest .

# 或者使用 buildx 多平台构建
docker buildx build --platform linux/amd64,linux/arm64 -t k8doctor:latest .
```

### 2. 部署到集群

```bash
# 方式一：交互式诊断（推荐）
kubectl apply -f k8s-deployment.yaml
kubectl exec -it -n k8doctor k8doctor -- /app/k8doctor diagnose-i

# 方式二：执行单次诊断
kubectl apply -f k8s-job.yaml

# 查看诊断结果
kubectl logs -n k8doctor job/k8doctor-diagnose
```

### 3. 使用部署脚本

```bash
# Linux/Mac
chmod +x deploy.sh
./deploy.sh build    # 构建镜像
./deploy.sh deploy   # 部署
./deploy.sh exec     # 进入交互式诊断

# Windows
deploy.bat build
deploy.bat deploy
deploy.bat exec
```

---

## 部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Target K8s Cluster                       │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              k8doctor Namespace                       │   │
│  │                                                        │   │
│  │  ┌────────────────────────────────────────────────┐   │   │
│  │  │  k8doctor Pod                                  │   │   │
│  │  │  ┌──────────────────────────────────────────┐ │   │   │
│  │  │  │ k8doctor Container                       │ │   │   │
│  │  │  │ - Go CLI (诊断工具)                       │ │   │   │
│  │  │  │ - /app/kb (知识库, 17个故障场景)         │ │   │   │
│  │  │  │ - /app/scripts (ONNX 推理脚本)           │ │   │   │
│  │  │  └──────────────────────────────────────────┘ │   │   │
│  │  └────────────────────────────────────────────────┘   │   │
│  │                                                        │   │
│  │  ServiceAccount: k8doctor                              │   │
│  │  ClusterRoleBinding: 只读访问集群资源                   │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 权限说明

部署时会创建 `k8doctor-reader` ClusterRole，授予以下只读权限：

| API 组 | 资源 | 权限 |
|--------|------|------|
| core | pods, pods/log | GET, LIST, WATCH |
| core | services, endpoints | GET, LIST, WATCH |
| core | nodes, events | GET, LIST, WATCH |
| core | configmaps, secrets | GET, LIST, WATCH |
| apps | deployments, replicasets | GET, LIST, WATCH |
| networking | ingresses | GET, LIST, WATCH |

---

## 使用场景

### 场景一：交互式诊断

```bash
# 部署交互式诊断环境
kubectl apply -f k8s-deployment.yaml

# 进入容器
kubectl exec -it -n k8doctor k8doctor -- /bin/bash

# 在容器内执行诊断
/app/k8doctor diagnose-i
# 或
/app/k8doctor diagnose "Pod一直重启，日志报OOM"
```

### 场景二：单次 Job 诊断

```bash
# 设置环境变量
export DIAGNOSE_SYMPTOM="Pod一直重启，日志报OOM"

# 或直接执行
DIAGNOSE_SYMPTOM="Ingress 502错误" kubectl apply -f k8s-job.yaml

# 查看结果
kubectl logs -n k8doctor job/k8doctor-diagnose
```

### 场景三：集成到 CI/CD

```yaml
# 在 CI/CD Pipeline 中使用
- name: K8s Health Check
  image: k8doctor:latest
  env:
    - name: DIAGNOSE_SYMPTOM
      value: "deployment progress not ready"
  command: ["/bin/sh", "-c"]
  args:
    - |
      /app/k8doctor diagnose "$DIAGNOSE_SYMPTOM"
```

---

## 配置说明

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `K8DOCTOR_KB_PATH` | `/app/kb` | 知识库路径 |
| `K8DOCTOR_EMBEDDINGS_PATH` | `/app/embeddings` | 向量索引路径 |

### 挂载卷

| 卷 | 类型 | 说明 |
|----|------|------|
| `kb-volume` | ConfigMap | 只读知识库 |
| `embeddings-volume` | EmptyDir | 向量索引（可写） |

---

## 镜像优化

当前镜像约 200MB（Alpine + Go binary）。如需进一步优化：

### 使用 distroless 镜像

```dockerfile
FROM gcr.io/distroless/static:nonroot AS runtime
COPY --from=builder /build/k8doctor /k8doctor
ENTRYPOINT ["/k8doctor"]
```

### 使用 scratch 镜像（需静态编译）

```dockerfile
FROM scratch AS runtime
COPY --from=builder /build/k8doctor /k8doctor
COPY kb /kb
ENTRYPOINT ["/k8doctor"]
```

---

## 知识库更新

更新知识库后需要重新构建镜像：

```bash
# 1. 修改 kb/*.md 文件

# 2. 重新构建镜像
docker build -t k8doctor:latest .

# 3. 更新部署
kubectl rollout restart deployment/k8doctor -n k8doctor
```

---

## 故障排除

### Pod 无法启动

```bash
# 查看事件
kubectl describe pod k8doctor -n k8doctor

# 查看日志
kubectl logs k8doctor -n k8doctor
```

### 权限不足

```bash
# 检查 ClusterRoleBinding
kubectl get clusterrolebinding k8doctor-reader-binding

# 手动创建权限（如果自动创建失败）
kubectl apply -f k8s-deployment.yaml
```

### 镜像拉取失败

```bash
# 使用本地镜像（适用于没有镜像仓库的环境）
kubectl apply -f k8s-deployment.yaml

# 或者将镜像导入到集群的镜像仓库
docker tag k8doctor:latest <your-registry>/k8doctor:latest
docker push <your-registry>/k8doctor:latest
```

---

## 资源限制

| 资源 | 请求 | 限制 |
|------|------|------|
| Memory | 128Mi | 512Mi |
| CPU | 100m | 500m |

生产环境可根据实际需求调整。
