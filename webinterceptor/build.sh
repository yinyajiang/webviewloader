#!/bin/bash

# 判断是否存在虚拟环境
if [ ! -f "pyinstaller_venv/bin/python3" ]; then
    # 创建虚拟环境
    python3 -m venv pyinstaller_venv
fi

# 激活虚拟环境
source pyinstaller_venv/bin/activate

python3 -V

pip3 install -r requirements.txt
pip3 install pyinstaller 
pip3 install requests

# 构建
rm -rf build dist .eggs
python3 build_pyinstaller.py --must-cert --onedir $@
cert=$(python3 build_cert.py)

# 获取dist的.app目录的名字
app_name=$(ls dist | grep -E '\.app$')

# Check if app was found
if [ -z "$app_name" ]; then
    echo "Error: No .app file found in dist directory"
    exit 1
fi

# 压缩zip
cd dist && zip -ry $app_name.zip $app_name && cd ..
# 计算zip的md5 (使用macOS的md5命令)
echo -n "$app_name.zip: " > dist/$app_name.zip.md5
md5 -q dist/$app_name.zip >> dist/$app_name.zip.md5
python3 -V
echo $cert

