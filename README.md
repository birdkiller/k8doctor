# K8s Doctor - Kubernetes 智能故障诊断工具

基于知识库的 K8s 故障诊断工具，通过语义理解匹配故障场景，给出完整可行的排查建议。

## 功能特性

- 🎯 **语义诊断**：输入不规则的故障症状描述，程序自动清洗并匹配知识库
- 📚 **知识库驱动**：Markdown 格式知识库，支持增量添加，适配 AI 记忆
- 🔍 **多维匹配**：Embedding 语义匹配 + 关键词兜底，最大化故障识别覆盖率
- 📋 **完整排查**：返回结构化排查步骤和可执行 kubectl 命令
- 🤖 **真语义向量**：支持 all-MiniLM-L6-v2 ONNX 模型（需单独安装）

## 支持的故障场景

| 类别 | 场景 |
|------|------|
| Pod | OOMKilled、CrashLoopBackOff、ImagePullBackOff、Pending、NotReady、Evicted |
| Node | NotReady、MemoryPressure、DiskPressure |
| 网络 | Service 不通、Endpoints 为空、DNS 解析失败、Ingress 502/504 |
| 存储 | PVC Pending、存储挂载失败 |
| RBAC | 权限拒绝、Secret 读取失败 |
| 部署 | Deployment 卡在 Progressing |

## 快速开始

### 1. 构建

```bash
cd D:\workspace\k8doctor
go mod tidy
go build -o k8doctor.exe ./cmd/cli
```

### 2. 运行

```bash
# 诊断（使用 TF-IDF 向量引擎）
./k8doctor.exe diagnose "Pod一直重启，日志报OOM killed"

# 交互式诊断
./k8doctor.exe diagnose-i

# 列出支持的故障类型
./k8doctor.exe list

# 重建向量索引
./k8doctor.exe rebuild-index
```

## 项目结构

```
k8doctor/
├── cmd/cli/main.go           # CLI 入口
├── internal/
│   ├── cleaner/             # 症状清洗模块
│   ├── knowledge/           # 知识库加载模块
│   ├── matcher/             # 向量匹配模块
│   └── output/              # 结果输出模块
├── kb/                      # 知识库 Markdown 文件 (17个场景)
├── embeddings/
│   ├── models/              # ONNX 模型 (需单独安装)
│   │   └── all-MiniLM-L6-v2.onnx
│   └── index.json           # 向量索引
├── scripts/
│   ├── download_model.py    # 模型下载脚本
│   ├── onnx_inference.py   # ONNX 推理脚本
│   ├── test_onnx.py        # 测试脚本
│   ├── requirements.txt    # Python 依赖
│   └── setup_onnx_manual.md # ONNX 安装指南
└── README.md
```

## 向量引擎

程序内置双轨向量引擎：

| 引擎 | 说明 | 状态 |
|------|------|------|
| **TF-IDF Fallback** | 纯 Go 实现，无需外部依赖 | ✅ 默认启用 |
| **all-MiniLM-L6-v2** | Google 预训练语义向量模型，384维 | 🔧 需安装 |

### 启用 ONNX 模型

ONNX 模型提供更精准的语义理解能力。如需启用：

1. 安装 Python 依赖：
   ```bash
   pip install numpy onnxruntime
   ```

2. 查看详细安装指南：
   ```
   scripts/setup_onnx_manual.md
   ```

3. 测试 ONNX：
   ```bash
   python scripts/test_onnx.py
   ```

程序会自动检测 ONNX 是否可用。

## 技术架构

```
┌──────────────────────────────────────────────────────────┐
│                      Go CLI                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │  cleaner    │→ │   matcher    │→ │   output     │     │
│  │  症状清洗   │  │  向量匹配   │  │  结果输出   │     │
│  └─────────────┘  └──────┬──────┘  └─────────────┘     │
│                          │                              │
│                    ┌─────▼─────┐                        │
│                    │ Python 子进程 │ ← ONNX 推理       │
│                    └───────────┘                        │
└──────────────────────────────────────────────────────────┘
```

## 添加新故障场景

在 `kb/` 目录下创建新的 `.md` 文件：

```markdown
---
id: pod-mycase
category: pod
severity: warning
title: "My Custom Case"
related_states: [MyCase, mycase]
tags: [custom, pod]
---

## 故障现象
描述故障现象...

## 排查步骤
1. kubectl get pods -n <namespace>
2. kubectl describe pod <pod-name> -n <namespace>

## 可执行命令
kubectl get pods -n <namespace>
```

然后重建索引：
```bash
./k8doctor.exe rebuild-index
```

## 目录内容

```
kb/
├── pod-oomkilled.md           # Pod OOMKilled
├── pod-crashloop.md          # Pod CrashLoopBackOff
├── pod-imagepullbackoff.md   # Pod ImagePullBackOff
├── pod-pending.md             # Pod Pending (调度失败)
├── pod-notready.md           # Pod Running 但 NotReady
├── pod-evicted.md             # Pod 被驱逐
├── node-notready.md           # Node NotReady
├── node-memorypressure.md     # Node 内存压力
├── node-diskpressure.md       # Node 磁盘压力
├── service-unreachable.md     # Service 无法访问
├── endpoints-empty.md          # Endpoints 为空
├── dns-failure.md             # DNS 解析失败
├── ingress-502.md             # Ingress 502/504 错误
├── pvc-pending.md             # PVC Pending
├── rbac-denied.md             # RBAC 权限拒绝
├── deployment-stuck.md        # Deployment 卡住
├── secret-read-failed.md      # Secret 读取失败
└── storage-mount-failed.md    # 存储挂载失败
```

## 容器化部署

支持 Docker 容器化部署到 K8s 集群进行故障诊断：

```bash
# 构建镜像
docker build -t k8doctor:latest .

# 部署到集群
kubectl apply -f k8s-deployment.yaml

# 进入交互式诊断
kubectl exec -it -n k8doctor k8doctor -- /app/k8doctor diagnose-i
```

详细文档参见 [DEPLOY_K8S.md](DEPLOY_K8S.md)

## License

MIT
