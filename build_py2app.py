from setuptools import setup
import argparse
import sys

parser = argparse.ArgumentParser(description='Command Line Parser')
parser.add_argument('--name', default='load_cookie')
parser.add_argument('--icon', default='')
parser.add_argument('*')
args = parser.parse_args()
for arg in sys.argv:
    if arg in ['--icon', '--name']:
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