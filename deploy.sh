#!/bin/bash
# K8s Doctor 部署脚本
# 用于构建 Docker 镜像并部署到 K8s 集群

set -e

# 配置
IMAGE_NAME="k8doctor"
IMAGE_TAG="latest"
NAMESPACE="k8doctor"
DEPLOYMENT_FILE="k8s-deployment.yaml"
JOB_FILE="k8s-job.yaml"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 is required but not installed."
        exit 1
    fi
}

# 显示帮助
show_help() {
    cat << EOF
K8s Doctor 部署脚本

用法: $0 [命令]

命令:
    build       构建 Docker 镜像
    deploy      部署到 K8s 集群
    job         以 Job 方式执行单次诊断
    exec        部署并进入交互式诊断 shell
    logs        查看诊断 pod 日志
    delete      删除部署
    help        显示此帮助信息

示例:
    $0 build              # 构建镜像
    $0 deploy             # 部署交互式诊断工具
    $0 exec               # 部署并进入 shell
    $0 job "Pod重启 OOM"  # 执行单次诊断

环境变量:
    DIAGNOSE_SYMPTOM      诊断症状（用于 job 模式）

EOF
}

# 构建 Docker 镜像
cmd_build() {
    log_info "Building Docker image: ${IMAGE_NAME}:${IMAGE_TAG}"
    
    # 检查 Docker
    check_command docker
    
    # 构建
    docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
    
    log_info "Docker image built successfully!"
    log_info "Image: ${IMAGE_NAME}:${IMAGE_TAG}"
}

# 部署到 K8s
cmd_deploy() {
    log_info "Deploying to Kubernetes namespace: ${NAMESPACE}"
    
    # 检查 kubectl
    check_command kubectl
    
    # 创建命名空间
    kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
    
    # 部署
    kubectl apply -f ${DEPLOYMENT_FILE}
    
    log_info "Deployment created!"
    log_info "Run 'kubectl exec -it -n ${NAMESPACE} k8doctor -- /app/k8doctor diagnose-i' to start diagnosing"
}

# 执行 Job 诊断
cmd_job() {
    SYMPTOM="${DIAGNOSE_SYMPTOM:-Pod一直重启，日志报OOM}"
    
    log_info "Running diagnosis job with symptom: ${SYMPTOM}"
    
    # 检查 kubectl
    check_command kubectl
    
    # 创建命名空间
    kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
    
    # 设置环境变量
    export DIAGNOSE_SYMPTOM="${SYMPTOM}"
    
    # 部署 Job
    kubectl apply -f ${JOB_FILE}
    
    # 等待 Job 完成
    log_info "Waiting for job to complete..."
    kubectl wait --for=condition=complete -n ${NAMESPACE} job/k8doctor-diagnose --timeout=60s
    
    # 查看日志
    log_info "Job output:"
    kubectl logs -n ${NAMESPACE} job/k8doctor-diagnose
}

# 进入交互式 shell
cmd_exec() {
    log_info "Deploying and entering interactive shell..."
    
    # 检查 kubectl
    check_command kubectl
    
    # 创建命名空间
    kubectl create namespace ${NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
    
    # 部署
    kubectl apply -f ${DEPLOYMENT_FILE}
    
    # 等待 Pod 就绪
    log_info "Waiting for pod to be ready..."
    kubectl wait --for=condition=ready -n ${NAMESPACE} pod/k8doctor --timeout=60s
    
    # 进入 shell
    kubectl exec -it -n ${NAMESPACE} k8doctor -- /bin/bash
}

# 查看日志
cmd_logs() {
    kubectl logs -n ${NAMESPACE} k8doctor
}

# 删除部署
cmd_delete() {
    log_info "Deleting deployment from namespace: ${NAMESPACE}"
    kubectl delete -f ${DEPLOYMENT_FILE} --ignore-not-found=true
    kubectl delete namespace ${NAMESPACE} --ignore-not-found=true
    log_info "Cleanup complete"
}

# 主逻辑
case "${1:-help}" in
    build)
        cmd_build
        ;;
    deploy)
        cmd_deploy
        ;;
    job)
        cmd_job
        ;;
    exec)
        cmd_exec
        ;;
    logs)
        cmd_logs
        ;;
    delete)
        cmd_delete
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        log_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
