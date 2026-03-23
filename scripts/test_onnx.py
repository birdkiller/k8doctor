#!/usr/bin/env python3
"""
测试 ONNX 模型推理是否正常工作

用法:
    python test_onnx.py
"""

import sys
import os

def main():
    print("=" * 60)
    print("测试 all-MiniLM-L6-v2 ONNX 推理")
    print("=" * 60)
    
    script_dir = os.path.dirname(os.path.abspath(__file__))
    model_path = os.path.join(os.path.dirname(script_dir), "embeddings", "models", "all-MiniLM-L6-v2.onnx")
    
    # 1. 检查模型文件
    print(f"\n1. 检查模型文件...")
    if os.path.exists(model_path):
        size_mb = os.path.getsize(model_path) / (1024 * 1024)
        print(f"   ✓ 模型存在: {model_path}")
        print(f"   ✓ 文件大小: {size_mb:.1f} MB")
    else:
        print(f"   ✗ 模型不存在: {model_path}")
        print(f"\n正在下载模型...")
        os.system(f'"{sys.executable}" "{os.path.join(script_dir, "download_model.py")}"')
        
        if os.path.exists(model_path):
            size_mb = os.path.getsize(model_path) / (1024 * 1024)
            print(f"   ✓ 模型下载成功: {size_mb:.1f} MB")
        else:
            print("   ✗ 模型下载失败")
            sys.exit(1)
    
    # 2. 检查依赖
    print(f"\n2. 检查 Python 依赖...")
    try:
        import numpy as np
        print(f"   ✓ numpy {np.__version__}")
    except ImportError:
        print("   ✗ numpy 未安装")
        print("   请运行: pip install numpy onnxruntime")
        sys.exit(1)
    
    try:
        import onnxruntime as ort
        print(f"   ✓ onnxruntime {ort.__version__}")
    except ImportError:
        print("   ✗ onnxruntime 未安装")
        print("   请运行: pip install numpy onnxruntime")
        sys.exit(1)
    
    # 3. 测试推理
    print(f"\n3. 测试推理...")
    try:
        import onnxruntime as ort
        import numpy as np
        
        session = ort.InferenceSession(model_path, providers=['CPUExecutionProvider'])
        print(f"   ✓ 模型加载成功")
        
        # 准备输入
        input_name = session.get_inputs()[0].name
        test_text = "kubernetes pod oomkilled crash"
        
        # 简单分词
        vocab = {"[CLS]": 101, "[SEP]": 102, "[PAD]": 0, "[UNK]": 100}
        words = test_text.lower().split()
        for i, w in enumerate(words[:250]):
            vocab[w] = 1000 + i  # 简化
        
        tokens = [vocab.get("[CLS]", 101)]
        for w in words:
            tokens.append(vocab.get(w, vocab["[UNK]"]))
        tokens.append(vocab.get("[SEP]", 102))
        
        # Padding
        max_len = 256
        attention_mask = [1] * len(tokens) + [0] * (max_len - len(tokens))
        input_ids = tokens + [0] * (max_len - len(tokens))
        
        # 推理
        outputs = session.run(None, {
            input_name: [input_ids],
            "attention_mask": [attention_mask]
        })
        
        embedding = outputs[0][0]
        print(f"   ✓ 推理成功")
        print(f"   ✓ 输出维度: {len(embedding)}")
        print(f"   ✓ 前5个值: {embedding[:5]}")
        
    except Exception as e:
        print(f"   ✗ 推理失败: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
    
    print("\n" + "=" * 60)
    print("✓ 所有测试通过！ONNX 模型工作正常")
    print("=" * 60)

if __name__ == "__main__":
    main()
