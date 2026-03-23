#!/usr/bin/env python3
"""
设置 temp 目录并安装 onnxruntime
"""

import os
import sys
import tempfile

def main():
    print(f"Default temp directories:")
    print(f"  TEMP: {os.environ.get('TEMP', 'not set')}")
    print(f"  TMP: {os.environ.get('TMP', 'not set')}")
    print(f"  TMPDIR: {os.environ.get('TMPDIR', 'not set')}")
    
    # Try to create a temp directory in a writable location
    test_dirs = [
        "C:\\Users\\kingdee\\AppData\\Local\\Temp",
        "C:\\temp",
        "C:\\tmp",
        os.path.expanduser("~"),
        os.path.dirname(__file__),
    ]
    
    for d in test_dirs:
        try:
            os.makedirs(d, exist_ok=True)
            test_file = os.path.join(d, "test_write.tmp")
            with open(test_file, "w") as f:
                f.write("test")
            os.remove(test_file)
            print(f"  Writable: {d}")
        except Exception as e:
            print(f"  Not writable: {d} - {e}")
    
    # Try to set tempdir manually
    temp_dir = os.path.join(os.path.dirname(__file__), "temp")
    try:
        os.makedirs(temp_dir, exist_ok=True)
        tempfile.tempdir = temp_dir
        print(f"\nSet tempfile.tempdir to: {temp_dir}")
    except Exception as e:
        print(f"Could not set tempfile.tempdir: {e}")

if __name__ == "__main__":
    main()
