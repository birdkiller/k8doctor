# ONNX 模型手动安装指南

## 当前状态

由于沙箱环境权限限制，无法自动安装 `onnxruntime` Python 包。
**项目已内置 TF-IDF Fallback 引擎，可以正常使用。**

如需启用真正的 all-MiniLM-L6-v2 语义向量，请在有权限的环境中执行以下步骤：

---

## 步骤 1: 安装 Python 依赖

```bash
# 确保 Python 3.9+ 已安装
python --version

# 安装依赖
pip install numpy onnxruntime
```

如果 pip 遇到权限问题，使用 `--user` 标志：
```bash
pip install numpy onnxruntime --user
```

---

## 步骤 2: 下载 ONNX 模型

模型会在首次运行时自动下载，也可以手动下载：

```bash
# 方法1: 通过 Python 脚本（推荐）
python scripts/download_model.py

# 方法2: 手动下载
# 从 HuggingFace 下载:
# https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2/resolve/main/onnx/model.onnx
# 保存到: embeddings/models/all-MiniLM-L6-v2.onnx
```

---

## 步骤 3: 验证安装

```bash
# 测试 ONNX 是否正常工作
python scripts/test_onnx.py
```

预期输出：
```
============================================================
测试 all-MiniLM-L6-v2 ONNX 推理
============================================================

1. 检查模型文件...
   ✓ 模型存在: ...\embeddings\models\all-MiniLM-L6-v2.onnx
   ✓ 文件大小: 89.3 MB

2. 检查 Python 依赖...
   ✓ numpy 2.x.x
   ✓ onnxruntime 1.x.x

3. 测试推理...
   ✓ 模型加载成功
   ✓ 推理成功
   ✓ 输出维度: 384
   ✓ 前5个值: [0.123, -0.456, ...]

============================================================
✓ 所有测试通过！ONNX 模型工作正常
============================================================
```

---

## 步骤 4: 构建并运行 Go 程序

```bash
# 构建
go mod tidy
go build -o k8doctor.exe ./cmd/cli

# 首次运行会自动使用 ONNX 引擎
./k8doctor.exe diagnose "Pod一直重启，日志报OOM killed"
```

---

## 手动下载模型（备选方案）

如果网络下载有问题，可以手动下载：

1. 访问: https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2/tree/main/onnx
2. 下载 `model.onnx` 文件
3. 保存到项目目录: `D:\workspace\k8doctor\embeddings\models\all-MiniLM-L6-v2.onnx`

---

## TF-IDF Fallback 说明

如果 ONNX 不可用，程序会自动使用 TF-IDF 向量引擎：

- ✅ 完全正常工作
- ✅ 无需外部依赖
- ✅ 关键词匹配仍然精准
- ⚠️ 语义理解能力较弱（不如真正的 Transformer 模型）

对于 K8s 故障诊断场景，TF-IDF + 关键词双路匹配已经能覆盖大部分场景，ONNX 只是增强语义理解能力。

---

## 常见问题

### Q: pip install 失败？
A: 确保有写入权限，或使用 `--user` 标志。如果问题持续，可能需要联系系统管理员。

### Q: 模型下载很慢？
A: 可以使用镜像源或手动下载。参见上文"手动下载模型"部分。

### Q: 仍然导入失败？
A: 检查 Python 版本（需要 3.9+）和 pip 版本（尝试 `python -m pip install --upgrade pip`）。
