import os
import sys
import argparse
import requests
import shutil
from build_cert import get_cert

current_dir = os.path.dirname(os.path.abspath(__file__))
os.chdir(current_dir)
iswin = sys.platform.startswith('win')
print(f'\nsys.prefix: {sys.prefix}\n\n')

parser = argparse.ArgumentParser(description='Command Line Parser')
parser.add_argument('--onedir', action='store_true')
parser.add_argument('--name', default='WebInterceptor')
parser.add_argument('--must-cert', action='store_true')
parser.add_argument('--icon', default='')
args = parser.parse_args()

cert = get_cert() if args.must_cert else ""

if args.icon and args.icon.startswith('http'):
    # 下载
    dest = os.path.join(current_dir, "icon.ico")
    response = requests.get(args.icon)
    with open(dest, 'wb') as f:
        f.write(response.content)
    args.icon = dest

import PyInstaller.__main__

pyinstaller_args = [
    os.path.join(".", "main.py"),
    "-y",
    "--windowed",
    f"--name={args.name}",
] + (["--onedir"] if args.onedir else ["--onefile"]) + (["--codesign-identity", cert] if cert else []) + ([
   
] if not iswin else [
    
]) + ([
    f"--icon={args.icon}"
] if args.icon else [])

PyInstaller.__main__.run(pyinstaller_args)

_tmp = os.path.join(current_dir, "dist", args.name)
if os.path.exists(_tmp):
    shutil.rmtree(_tmp)
