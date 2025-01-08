import os
import sys
import argparse
import requests
import shutil
import subprocess
current_dir = os.path.dirname(os.path.abspath(__file__))
os.chdir(current_dir)
iswin = sys.platform.startswith('win')
print(f'\nsys.prefix: {sys.prefix}\n\n')

parser = argparse.ArgumentParser(description='Command Line Parser')
parser.add_argument('--onedir', action='store_true')
parser.add_argument('--name', default='load_cookie')
parser.add_argument('--cert', default='')
parser.add_argument('--icon', default='')
parser.add_argument('--win-sign', default='')
parser.add_argument('--bundle-id', default='')
args = parser.parse_args()

if args.bundle_id:
    raise Exception('not implement bundleid')


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
    "--exclude-module=PyQt5",
    "--exclude-module=PyQt6",
    "--exclude-module=PySide2",
    "--exclude-module=PySide6",
    f"--name={args.name}",
] + (["--onedir"] if args.onedir else ["--onefile"]) + (["--codesign-identity", args.cert, "--no-entitlements"] if args.cert else []) + ([
    "--hidden-import=WebKit",
    "--hidden-import=Foundation",
    "--hidden-import=webview",
] if not iswin else [
    "--collect-binaries=clr_loader",
]) + ([
    f"--icon={args.icon}"
] if args.icon else [])

PyInstaller.__main__.run(pyinstaller_args)

_tmp = os.path.join(current_dir, "dist", args.name)
if not args.onedir and os.path.exists(_tmp):
    shutil.rmtree(_tmp)

if iswin and args.win_sign:
    subprocess.run(args.win_sign, shell=True)

