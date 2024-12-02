import os
import sys
import argparse
from build_cert import get_cert


os.chdir(os.path.dirname(os.path.abspath(__file__)))
iswin = sys.platform.startswith('win')
print(f'\nsys.prefix: {sys.prefix}\n\n')

parser = argparse.ArgumentParser(description='Command Line Parser')
parser.add_argument('--onedir', action='store_true')
parser.add_argument('--name', default='load_cookie')
parser.add_argument('--must-cert', action='store_true')
parser.add_argument('--icon', default='')
args = parser.parse_args()

cert = get_cert() if args.must_cert else ""

import PyInstaller.__main__

pyinstaller_args = [
    os.path.join(".", "webview.py"),
    "-y",
    "--exclude-module=PyQt5",
    "--exclude-module=PyQt6",
    "--exclude-module=PySide2",
    "--exclude-module=PySide6",
    f"--name={args.name}",
] + (["--onedir"] if args.onedir else ["--onefile"]) + (["--codesign-identity", cert, "--no-entitlements"] if cert else []) + ([
    "--hidden-import=WebKit",
    "--hidden-import=Foundation",
    "--hidden-import=webview",
] if not iswin else []) + ([
    f"--icon={args.icon}"
] if args.icon else [])

PyInstaller.__main__.run(pyinstaller_args)
