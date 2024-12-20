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

#解析传过来的cert
while [[ $# -gt 0 ]]; do
    case $1 in
        --cert=*)
            cert="${1#*=}"
            shift
            ;;
        --cert)
            cert="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done
echo "cert: $cert"

# 构建
rm -rf build dist .eggs
python3 build_py2app.py py2app $@

# 获取dist的.app目录的名字
app_name=$(ls dist | grep -E '\.app$')

# Check if app was found
if [ -z "$app_name" ]; then
    echo "Error: No .app file found in dist directory"
    exit 1
fi

# 只有当cert有值时才进行签名
if [ ! -z "$cert" ]; then
    echo "codesign $app_name ... ..."
    if codesign --deep --force --verify --verbose --timestamp --sign "$cert" dist/$app_name; then
        echo "codesign $app_name success"
    else
        echo "codesign failed"
        exit 1
    fi
fi

# 压缩zip
cd dist && zip -r $app_name.zip $app_name && cd ..

# 计算zip的md5 (使用macOS的md5命令)
echo -n "$app_name.zip: " > dist/$app_name.zip.md5
md5 -q dist/$app_name.zip >> dist/$app_name.zip.md5
python3 -V
