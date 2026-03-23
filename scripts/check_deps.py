import sys
print(f"Python version: {sys.version}")

try:
    import numpy
    print(f"numpy: {numpy.__version__}")
except ImportError:
    print("numpy: NOT INSTALLED")

try:
    import onnxruntime
    print(f"onnxruntime: {onnxruntime.__version__}")
except ImportError:
    print("onnxruntime: NOT INSTALLED")
