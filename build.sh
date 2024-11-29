#!/bin/bash

# 判断是否存在虚拟环境
if [ ! -d "venv" ]; then
    # 创建虚拟环境
    python3 -m venv venv
fi

# 激活虚拟环境
source venv/bin/activate

pip3 install -r requirements.txt
pip3 install pyinstaller

# 构建
python3 pyinstaller_build.py $@
