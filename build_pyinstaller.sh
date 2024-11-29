#!/bin/bash

# 判断是否存在虚拟环境
if [ ! -d "pyinstaller_venv" ]; then
    # 创建虚拟环境
    python3 -m venv pyinstaller_venv
fi

# 激活虚拟环境
source pyinstaller_venv/bin/activate

pip3 install -r requirements.txt
pip3 install pyinstaller

# 构建
python3 build_pyinstaller.py $@
