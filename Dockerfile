# K8s Doctor - Kubernetes 智能故障诊断工具
# 用于容器化部署到 K8s 集群进行故障诊断

# ============================================================
# Stage 1: Builder
# ============================================================
FROM golang:1.22-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖（利用 Docker 缓存）
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
# CGO_ENABLED=0 是因为我们使用纯 Go 构建，不需要 cgo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o k8doctor \
    ./cmd/cli

# ============================================================
# Stage 2: ONNX 模型下载器 (可选)
# ============================================================
FROM python:3.12-slim AS model-downloader

WORKDIR /download

COPY scripts/download_model.py .

RUN pip install --no-cache-dir numpy onnxruntime && \
    python download_model.py

# ============================================================
# Stage 3: 运行镜像
# ============================================================
FROM alpine:3.19 AS runtime

# 安装运行时依赖
RUN apk add --no-cache \
    ca-certificates \
    bash \
    curl \
    kubectl \
    && update-ca-certificates

# 创建工作目录
WORKDIR /app

# 从 builder 复制二进制文件
COPY --from=builder /build/k8doctor /app/k8doctor

# 复制知识库
COPY --from=builder /build/kb /app/kb

# 复制脚本（包含 ONNX 推理脚本）
COPY --from=builder /build/scripts /app/scripts

# 复制 ONNX 模型（如果有）
COPY --from=model-downloader /download/embeddings/models/*.onnx /app/embeddings/models/ 2>/dev/null || true

# 设置环境变量
ENV K8DOCTOR_KB_PATH=/app/kb
ENV K8DOCTOR_EMBEDDINGS_PATH=/app/embeddings
ENV PATH=/app:$PATH

# 创建符号链接，方便使用
RUN chmod +x /app/k8doctor && \
    ln -sf /app/k8doctor /usr/local/bin/k8doctor

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD /app/k8doctor list > /dev/null 2>&1 || exit 1

# 默认入口点
ENTRYPOINT ["/app/k8doctor"]

# 默认命令：交互式诊断
CMD ["diagnose-i"]
