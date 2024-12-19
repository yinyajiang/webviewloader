import os
import sys
import argparse
import requests
import shutil
from build_cert import get_cert
import subprocess
current_dir = os.path.dirname(os.path.abspath(__file__))
os.chdir(current_dir)
iswin = sys.platform.startswith('win')
print(f'\nsys.prefix: {sys.prefix}\n\n')

parser = argparse.ArgumentParser(description='Command Line Parser')
parser.add_argument('--onedir', action='store_true')
parser.add_argument('--name', default='WebInterceptor')
parser.add_argument('--must-cert', action='store_true')
parser.add_argument('--icon', default='')
parser.add_argument('--win-sign', default='')
args = parser.parse_args()
print(args)

cert = get_cert() if args.must_cert else ""

if args.icon and args.icon.startswith('http'):
    # 下载
    dest = os.path.join(current_dir, "icon.ico")
    response = requests.get(args.icon)
    with open(dest, 'wb') as f:
        f.write(response.content)
    args.icon = dest
elif args.icon:
    dest = os.path.join(current_dir, "icon.ico")
    shutil.copy2(args.icon, dest) 
    args.icon = dest
    
# 确保图标文件路径是相对于当前目录的
if args.icon:
    args.icon = os.path.abspath(args.icon)

import PyInstaller.__main__

pyinstaller_args = [
    os.path.join(".", "main.py"),
    "-y",
    "--windowed",
    f"--name={args.name}",
] + (["--onedir"] if args.onedir else ["--onefile"]) + (["--codesign-identity", cert] if cert else []) + ([
   
] if iswin else [
    "--target-architecture=universal2"
]) + ([f"--icon={args.icon}"] if args.icon else []) + ([f"--add-data={args.icon};."] if args.icon and iswin else [])

PyInstaller.__main__.run(pyinstaller_args)

_tmp = os.path.join(current_dir, "dist", args.name)
if (not args.onedir and os.path.exists(_tmp)) or (args.onedir and os.path.exists(_tmp) and not iswin):
    shutil.rmtree(_tmp)

if iswin and args.win_sign:
    # 运行签名命令，捕获输出并等待完成
    subprocess.run(args.win_sign, shell=True)
    print("sign done")

