from setuptools import setup
import argparse
import sys
import requests
import os

current_dir = os.path.dirname(os.path.abspath(__file__))
parser = argparse.ArgumentParser(description='Command Line Parser')
parser.add_argument('--name', default='load_cookie')
parser.add_argument('--icon', default='')
parser.add_argument('--cert', default='')
parser.add_argument('--bundle-id', default='com.example.webview')
parser.add_argument('*')
args = parser.parse_args()
if args.icon and args.icon.startswith('http'):
    dest = os.path.join(current_dir, "icon.ico")
    response = requests.get(args.icon)
    with open(dest, 'wb') as f:
        f.write(response.content)
    args.icon = dest

# 删除参数，否则失败
for arg in sys.argv:
    if arg in ['--icon', '--name', '--cert', '--bundle-id']:
        index = sys.argv.index(arg)
        sys.argv.pop(index)
        sys.argv.pop(index)

ENTRY_POINT = ['main.py']

DATA_FILES = []
OPTIONS = {
    'argv_emulation': False,
    'strip': True,
    'includes': ['WebKit', 'Foundation', 'webview'],
    'excludes': ['PyQt5', 'PyQt6', 'PySide2', 'PySide6'],
    'plist': {
        'CFBundleIdentifier': args.bundle_id.lower().replace(' ', ''), 
        'CFBundleShortVersionString': '1.0.0', 
        'CFBundleVersion': '1.0.0',
    }
}
if args.icon:
    OPTIONS['iconfile'] = 'icon.icns'

setup(
    name=f'{args.name}',
    app=ENTRY_POINT,
    data_files=DATA_FILES,
    options={'py2app': OPTIONS},
    setup_requires=['py2app'],
)