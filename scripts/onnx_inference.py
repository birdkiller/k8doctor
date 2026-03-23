#!/usr/bin/env python3
"""
all-MiniLM-L6-v2 ONNX 推理脚本
Go 代码通过调用此脚本获取文本 embedding

用法:
    python onnx_inference.py "<文本>"
    python onnx_inference.py --file <文本文件>
"""

import sys
import os
import json
import math

def cosine_normalize(vec):
    """L2 归一化"""
    norm = math.sqrt(sum(v * v for v in vec))
    if norm > 0:
        return [v / norm for v in vec]
    return vec

def main():
    if len(sys.argv) < 2:
        print("用法: python onnx_inference.py \"<文本>\"", file=sys.stderr)
        sys.exit(1)

    text = sys.argv[1]
    
    # 模型路径
    script_dir = os.path.dirname(os.path.abspath(__file__))
    model_path = os.path.join(os.path.dirname(script_dir), "embeddings", "models", "all-MiniLM-L6-v2.onnx")
    
    if not os.path.exists(model_path):
        # 模型不存在，尝试下载
        print("模型不存在，正在下载...", file=sys.stderr)
        import subprocess
        result = subprocess.run([sys.executable, os.path.join(script_dir, "download_model.py")])
        if result.returncode != 0 or not os.path.exists(model_path):
            print("模型下载失败", file=sys.stderr)
            sys.exit(1)

    try:
        import numpy as np
        import onnxruntime as ort
        
        # 加载模型
        session = ort.InferenceSession(model_path, providers=['CPUExecutionProvider'])
        
        # 获取模型信息
        inputs = session.get_inputs()
        input_name = inputs[0].name  # 通常是 "input_ids"
        
        # 简单的 BERT tokenizer
        vocab_path = os.path.join(os.path.dirname(script_dir), "embeddings", "models", "vocab.txt")
        tokens = tokenize(text, vocab_path)
        
        # Padding 到 256
        max_len = 256
        attention_mask = [1] * len(tokens) + [0] * (max_len - len(tokens))
        input_ids = tokens + [0] * (max_len - len(tokens))
        
        # 推理
        onnx_inputs = {
            input_name: [input_ids],
            "attention_mask": [attention_mask]
        }
        
        outputs = session.run(None, onnx_inputs)
        
        # all-MiniLM-L6-v2 输出是 pooled output, shape: [1, 384]
        embedding = outputs[0][0].tolist()
        
        # 归一化
        embedding = cosine_normalize(embedding)
        
        # 输出 JSON 格式
        result = {
            "embedding": embedding,
            "success": True
        }
        print(json.dumps(result))
        
    except ImportError as e:
        print(f"缺少依赖库: {e}", file=sys.stderr)
        print("请安装: pip install numpy onnxruntime", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"推理失败: {e}", file=sys.stderr)
        sys.exit(1)

def tokenize(text, vocab_path=None):
    """简单的 BERT tokenizer"""
    # 特殊 token
    vocab = {
        "[PAD]": 0,
        "[UNK]": 100,
        "[CLS]": 101,
        "[SEP]": 102,
        "[MASK]": 103,
    }
    
    # 加载词表
    if vocab_path and os.path.exists(vocab_path):
        with open(vocab_path, 'r', encoding='utf-8') as f:
            for i, line in enumerate(f):
                word = line.strip()
                if word and word not in vocab:
                    vocab[word] = len(vocab)
    
    # 基础词汇（当词表不存在时使用）
    if len(vocab) < 100:
        words = ["the", "a", "an", "is", "are", "was", "be", "have", "has", "had",
                "do", "does", "did", "will", "would", "could", "should",
                "to", "of", "in", "for", "on", "with", "at", "by", "from",
                "and", "but", "or", "not", "this", "that", "it", "as",
                "pod", "pods", "node", "service", "deployment", "container",
                "error", "fail", "pending", "running", "crash", "restart",
                "memory", "cpu", "disk", "network", "timeout", "oom"]
        for i, w in enumerate(words):
            if w not in vocab:
                vocab[w] = len(vocab)
    
    # 分词
    text = text.lower().strip()
    
    # 简单分词
    chars = []
    for c in text:
        if c.isalnum():
            chars.append(c)
        else:
            chars.append(' ')
    
    words = ''.join(chars).split()
    
    # 转换为 token IDs
    cls_id = vocab.get("[CLS]", 101)
    sep_id = vocab.get("[SEP]", 102)
    unk_id = vocab.get("[UNK]", 100)
    
    token_ids = [cls_id]
    for word in words[:250]:  # 留一个位置给 SEP
        token_ids.append(vocab.get(word, unk_id))
    token_ids.append(sep_id)
    
    return token_ids

if __name__ == "__main__":
    main()
