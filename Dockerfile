# K8s Doctor - Kubernetes 智能故障诊断工具
# 用于容器化部署到 K8s 集群进行故障诊断

# ============================================================
# Stage 1: Builder
# ============================================================
FROM golang:1.22 AS builder

# 安装构建依赖 (Debian base)
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# 复制 go mod 文件
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o k8doctor \
    ./cmd/cli

# ============================================================
# Stage 2: ONNX 模型下载器 (可选)
# ============================================================
FROM python:3.12 AS model-downloader

WORKDIR /download

# 安装 Python 依赖
RUN pip install --no-cache-dir numpy onnxruntime

# 下载模型（如果失败不影响构建）
COPY scripts/download_model.py .
RUN python download_model.py || echo "ONNX model download failed, will use TF-IDF fallback"

# ============================================================
# Stage 3: 运行镜像
# ============================================================
FROM alpine:3.19 AS runtime

# 安装运行时依赖
RUN apk add --no-cache \
    ca-certificates \
    bash \
    curl \
    wget \
    && rm -rf /var/cache/apk/*

# 安装 kubectl（从官方二进制下载）
RUN wget -q -O /usr/local/bin/kubectl https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl \
    && chmod +x /usr/local/bin/kubectl

# 创建工作目录
WORKDIR /app

# 创建必要的目录
RUN mkdir -p /app/kb /app/embeddings/models /app/scripts

# 从 builder 复制二进制文件
COPY --from=builder /build/k8doctor /app/k8doctor

# 复制知识库
COPY --from=builder /build/kb/ /app/kb/

# 复制脚本
COPY --from=builder /build/scripts/ /app/scripts/

# 复制 ONNX 模型（如果有）
COPY --from=model-downloader /download/embeddings/models/ /app/embeddings/models/ 2>/dev/null || true

# 设置环境变量
ENV K8DOCTOR_KB_PATH=/app/kb
ENV K8DOCTOR_EMBEDDINGS_PATH=/app/embeddings
ENV PATH=/app:$PATH

# 创建符号链接
RUN chmod +x /app/k8doctor && \
    ln -sf /app/k8doctor /usr/local/bin/k8doctor

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD /app/k8doctor list > /dev/null 2>&1 || exit 1

# 默认入口点
ENTRYPOINT ["/app/k8doctor"]

# 默认命令：交互式诊断
CMD ["diagnose-i"]
