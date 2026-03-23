#!/usr/bin/env python3
"""
下载 all-MiniLM-L6-v2 ONNX 模型
模型来源: HuggingFace sentence-transformers
"""

import os
import sys

def download_file(url, dest_path):
    """下载文件"""
    import urllib.request
    
    print(f"下载中: {url}")
    print(f"目标: {dest_path}")
    
    # 确保目录存在
    os.makedirs(os.path.dirname(dest_path), exist_ok=True)
    
    def reporthook(block_num, block_size, total_size):
        downloaded = block_num * block_size
        percent = min(100, downloaded * 100 // total_size) if total_size > 0 else 0
        sys.stdout.write(f"\r下载进度: {percent}% ({downloaded // (1024*1024)}MB / {total_size // (1024*1024)}MB)")
        sys.stdout.flush()
    
    urllib.request.urlretrieve(url, dest_path, reporthook)
    print(f"\n下载完成: {dest_path}")

def main():
    # 模型下载地址
    model_url = "https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2/resolve/main/onnx/model.onnx"
    backup_url = "https://hf-mirror.com/sentence-transformers/all-MiniLM-L6-v2/resolve/main/onnx/model.onnx"
    
    # 目标路径
    script_dir = os.path.dirname(os.path.abspath(__file__))
    models_dir = os.path.join(os.path.dirname(script_dir), "embeddings", "models")
    model_path = os.path.join(models_dir, "all-MiniLM-L6-v2.onnx")
    
    # 检查是否已存在
    if os.path.exists(model_path):
        size_mb = os.path.getsize(model_path) / (1024 * 1024)
        print(f"模型已存在: {model_path} ({size_mb:.1f}MB)")
        return
    
    # 尝试下载
    try:
        download_file(model_url, model_path)
    except Exception as e:
        print(f"\n主源下载失败: {e}")
        print("尝试备用源...")
        try:
            download_file(backup_url, model_path)
        except Exception as e2:
            print(f"备用源也失败: {e2}")
            sys.exit(1)

if __name__ == "__main__":
    main()
