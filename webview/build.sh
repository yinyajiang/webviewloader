#!/bin/bash

# 判断是否存在虚拟环境
if [ ! -f "py2app_venv/bin/python3" ]; then
    # 创建虚拟环境
    python3 -m venv py2app_venv
fi

# 激活虚拟环境
source py2app_venv/bin/activate

python3 -V

pip3 install setuptools==70.3.0
pip3 install -r requirements.txt
pip3 install py2app 
# 构建
rm -rf build dist .eggs
python3 build_py2app.py py2app $@
cert=$(python3 build_cert.py)

# 获取dist的.app目录的名字
app_name=$(ls dist | grep -E '\.app$')

# Check if app was found
if [ -z "$app_name" ]; then
    echo "Error: No .app file found in dist directory"
    exit 1
fi

echo "codesign $app_name ... ..."

if codesign --deep --force --verify --verbose --timestamp --sign "$cert" dist/$app_name; then
    echo "codesign $app_name success"

    # 压缩zip
    cd dist && zip -r $app_name.zip $app_name && cd ..

    # 计算zip的md5 (使用macOS的md5命令)
    echo -n "$app_name.zip: " > dist/$app_name.zip.md5
    md5 -q dist/$app_name.zip >> dist/$app_name.zip.md5
fi
