#!/usr/bin/env python3
"""
使用 Python 直接下载和安装 onnxruntime
"""

import subprocess
import sys
import os
import urllib.request
import zipfile
import shutil

def main():
    print("Attempting to install onnxruntime...")
    
    # Try direct pip install with different approach
    cmd = [sys.executable, "-m", "pip", "install", "onnxruntime"]
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=300)
        print("stdout:", result.stdout)
        print("stderr:", result.stderr)
        print("returncode:", result.returncode)
    except subprocess.TimeoutExpired:
        print("Install timed out")
    except Exception as e:
        print(f"Install exception: {e}")
    
    # Verify
    try:
        import onnxruntime
        print(f"SUCCESS: onnxruntime {onnxruntime.__version__}")
    except ImportError:
        print("FAILED: onnxruntime not installed")
        print("\nAlternative: Download wheel manually from:")
        print("https://pypi.org/project/onnxruntime/#files")

if __name__ == "__main__":
    main()
