#!/usr/bin/env python3
"""
手动下载和安装 onnxruntime wheel 文件
绕过 pip 的 temp 目录问题
"""

import os
import sys
import urllib.request
import zipfile
import site

def main():
    print("Manual onnxruntime installation")
    
    # Get Python paths
    user_site = site.getusersitepackages()
    print(f"User site packages: {user_site}")
    
    # Ensure directory exists
    os.makedirs(user_site, exist_ok=True)
    
    # onnxruntime wheel URL for Python 3.12 on Windows
    # We need to find the correct version for the platform
    wheel_url = "https://files.pythonhosted.org/packages/py2.py3/o/onnxruntime/onnxruntime-1.19.2-py3-none-any.whl"
    
    # Actually, onnxruntime has platform-specific wheels
    # Let's try a different approach - use the any wheel first
    # If that doesn't work, we'll need to find the correct platform wheel
    
    wheel_path = os.path.join(os.path.dirname(__file__), "onnxruntime.whl")
    
    if os.path.exists(wheel_path):
        print(f"Wheel already exists at {wheel_path}")
    else:
        print(f"Downloading onnxruntime wheel...")
        try:
            urllib.request.urlretrieve(wheel_url, wheel_path)
            print(f"Downloaded to {wheel_path}")
        except Exception as e:
            print(f"Download failed: {e}")
            print("\nAlternative approach: Use a lighter alternative")
            print("onnxruntime-slim or just skip ONNX for now")
            return False
    
    # Extract wheel
    print("Extracting wheel...")
    try:
        with zipfile.ZipFile(wheel_path, 'r') as zip_ref:
            zip_ref.extractall(user_site)
        print(f"Extracted to {user_site}")
    except Exception as e:
        print(f"Extraction failed: {e}")
        return False
    
    # Verify
    try:
        import onnxruntime
        print(f"SUCCESS: onnxruntime {onnxruntime.__version__}")
        return True
    except ImportError as e:
        print(f"Import failed: {e}")
        return False

if __name__ == "__main__":
    success = main()
    sys.exit(0 if success else 1)
