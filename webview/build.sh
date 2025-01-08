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


# 使用数组保存参数，这样可以正确处理带空格的值
declare -a original_args=()
cert=""

# 遍历所有参数
while [[ $# -gt 0 ]]; do
    if [[ "$1" == --cert=* ]]; then
        cert="${1#*=}"
    elif [[ "$1" == "--cert" && -n "$2" ]]; then
        cert="$2"
        shift  # 额外移动一次，跳过证书值
    else
        original_args+=("$1")
    fi
    shift
done

echo "cert: $cert"

# 构建
rm -rf build dist .eggs
python3 build_py2app.py py2app "${original_args[@]}"

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
