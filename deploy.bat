@echo off
setlocal EnableDelayedExpansion

set IMAGE_NAME=k8doctor
set IMAGE_TAG=latest
set NAMESPACE=k8doctor

echo ========================================
echo K8s Doctor 部署脚本
echo ========================================
echo.

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="/?" goto help
if "%1"=="-h" goto help

if "%1"=="build" goto build
if "%1"=="deploy" goto deploy
if "%1"=="job" goto job
if "%1"=="exec" goto exec
if "%1"=="logs" goto logs
if "%1"=="delete" goto delete

:help
echo 用法: %0 [命令]
echo.
echo 命令:
echo   build       构建 Docker 镜像
echo   deploy      部署到 K8s 集群
echo   job         以 Job 方式执行单次诊断
echo   exec        部署并进入交互式诊断 shell
echo   logs        查看诊断 pod 日志
echo   delete      删除部署
echo   help        显示此帮助信息
echo.
echo 示例:
echo   %0 build
echo   %0 deploy
echo   %0 job "Pod重启 OOM"
goto end

:build
echo [INFO] Building Docker image: %IMAGE_NAME%:%IMAGE_TAG%
docker build -t %IMAGE_NAME%:%IMAGE_TAG% .
echo [INFO] Docker image built successfully!
goto end

:deploy
echo [INFO] Deploying to Kubernetes namespace: %NAMESPACE%
kubectl create namespace %NAMESPACE% --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f k8s-deployment.yaml
echo [INFO] Deployment created!
echo [INFO] Run: kubectl exec -it -n %NAMESPACE% k8doctor -- /app/k8doctor diagnose-i
goto end

:job
set SYMPTOM=%DIAGNOSE_SYMPTOM%
if "%SYMPTOM%"=="" set SYMPTOM=Pod一直重启，日志报OOM
echo [INFO] Running diagnosis job with symptom: %SYMPTOM%
kubectl create namespace %NAMESPACE% --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f k8s-job.yaml
echo [INFO] Waiting for job to complete...
kubectl wait --for=condition=complete -n %NAMESPACE% job/k8doctor-diagnose --timeout=60s
echo [INFO] Job output:
kubectl logs -n %NAMESPACE% job/k8doctor-diagnose
goto end

:exec
echo [INFO] Deploying and entering interactive shell...
kubectl create namespace %NAMESPACE% --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f k8s-deployment.yaml
echo [INFO] Waiting for pod to be ready...
kubectl wait --for=condition=ready -n %NAMESPACE% pod/k8doctor --timeout=60s
kubectl exec -it -n %NAMESPACE% k8doctor -- /bin/bash
goto end

:logs
kubectl logs -n %NAMESPACE% k8doctor
goto end

:delete
echo [INFO] Deleting deployment from namespace: %NAMESPACE%
kubectl delete -f k8s-deployment.yaml --ignore-not-found=true
kubectl delete namespace %NAMESPACE% --ignore-not-found=true
echo [INFO] Cleanup complete
goto end

:end
