import os
import sys
import argparse
import subprocess
import re


os.chdir(os.path.dirname(os.path.abspath(__file__)))
iswin = sys.platform.startswith('win')
print(f'\nsys.prefix: {sys.prefix}\n\n')

parser = argparse.ArgumentParser(description='Command Line Parser')
parser.add_argument('--onedir', action='store_true')
parser.add_argument('--name', default='load_cookie')
parser.add_argument('--must-cert', action='store_true')
args = parser.parse_args()

cert = ""
if args.must_cert:
    prename='Developer ID Application:'
    output = subprocess.check_output(f'security find-certificate -c "{prename}"', shell=True).decode('utf-8')
    match = re.compile(f'"({prename}.+)"').search(output)
    cert = match.group(1)
    if cert == '':
        raise Exception('No certificate found')

import PyInstaller.__main__

pyinstaller_args = [
    os.path.join(".", "main.py"),
    "-y",
    "--exclude-module=PyQt5",
    "--exclude-module=PyQt6",
    "--exclude-module=PySide2",
    "--exclude-module=PySide6",
    f"--name={args.name}",
] + (["--onedir"] if args.onedir else ["--onefile"]) + (["--codesign-identity", cert] if cert else []) +([
    "--hidden-import=WebKit",
    "--hidden-import=Foundation",
    "--hidden-import=webview",
] if not iswin else [])

PyInstaller.__main__.run(pyinstaller_args)
