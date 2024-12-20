#!/usr/bin/env python3
import os
import sys
import shutil
import argparse
import subprocess
from pathlib import Path    
import re
import subprocess
import requests
import hashlib


def check_cert_valid(cert): 
    prename='Developer ID Application:'
    output = subprocess.check_output(f'security find-certificate -c "{prename}"', shell=True).decode('utf-8')
    if cert not in output:
        raise Exception(f'{cert} not found')
    return True


def getBatEnv(bat):
    p = subprocess.Popen('cmd /c "{}" ^&^& set'.format(bat),
                         shell=True,
                         stdout=subprocess.PIPE)
    outs, errs = p.communicate()
    env = {}
    for line in str(outs, encoding='utf-8').splitlines():
        s = line.split("=", 1)
        if len(s) != 2:
            continue
        k, v = s
        if k.lower() == "path":
            k = "PATH"
        env[k] = v
    return env


def setWinQtEnv(vsbat):
    envAll = getBatEnv(vsbat)
    path = set()
    for p in envAll["PATH"].split(";"):
        if p:
            path.add(p)

    for k, v in envAll.items():
        os.environ[k] = v


current_dir = Path(__file__).parent
os.chdir(current_dir)

def main():

    parser = argparse.ArgumentParser(description='Build Qt project')
    parser.add_argument('--name', default='webinterceptor',
                      help='Target executable name')
    parser.add_argument('--bundle-id', default='com.example.webinterceptor',
                      help='Bundle identifier (macOS)')
    parser.add_argument('--icon',
                      help='Path to icon file')
    parser.add_argument('--cert',
                      help='maccert')
    parser.add_argument('--win-sign',
                      help='sign cmd')
    parser.add_argument('--win-vsbat',
                      help='Path to win vsbat file')
    parser.add_argument('--qt-bin',
                      help='Path to qt bin')
    args = parser.parse_args()

    # 删除dist
    if os.path.exists('dist'):
        shutil.rmtree('dist')

    if args.icon and args.icon.startswith('http'):
        dest = os.path.join(current_dir, "dist", args.icon.split('/')[-1])
        response = requests.get(args.icon)
        with open(dest, 'wb') as f:
            f.write(response.content)
        args.icon = dest

    isWin= sys.platform == 'win32'

    if not isWin:
        with open('Info.plist', 'r') as file:
            content = file.read()
            if args.bundle_id:
                content = content.replace('{bundle_id}', args.bundle_id.replace(' ', '').lower())
        with open('Info_build.plist', 'w') as file:
            file.write(content)
    else:
        setWinQtEnv(args.win_vsbat)

    with open('webinterceptor.pro', 'r') as file:
        content = file.read()
        content = content.replace('{name}', args.name)
        if not isWin:
            content = content.replace('Info.plist', "Info_build.plist")
        if args.icon:
            if isWin:
                content = content.replace(';RC_ICONS', f'RC_ICONS = {args.icon}')
            else:
                content = content.replace(';ICON', f'ICON = {args.icon}')
    with open('webinterceptor_build.pro', 'w') as file:
        file.write(content)


    if not isWin:
        subprocess.run([f'{args.qt_bin}/qmake', 'webinterceptor_build.pro'], cwd=current_dir).check_returncode()
        subprocess.run(['make'], cwd=current_dir).check_returncode()
        subprocess.run([f'{args.qt_bin}/macdeployqt', f'dist/{args.name}.app',
                        ], cwd=current_dir).check_returncode()
    
        if args.cert:
            check_cert_valid(args.cert)
            subprocess.run([f'codesign', '--timestamp', '--force', '--deep', '--verify', '--verbose', '--sign', args.cert, f'dist/{args.name}.app'], cwd=current_dir).check_returncode()
            print(f'codesign success')
        else:
            print(f'not codesign')
        subprocess.run(['zip', '-ry', f'{args.name}.app.zip', f'{args.name}.app'], cwd=os.path.join(current_dir, 'dist')).check_returncode()
        print(f'zip success')
        with open(f"dist/{args.name}.app.zip.md5", 'w') as file:
            file.write(f"{args.name}.app.zip: " +  hashlib.md5(open(f"dist/{args.name}.app.zip", 'rb').read()).hexdigest())
    else:
        subprocess.run([f'{args.qt_bin}/qmake', 'webinterceptor_build.pro'], cwd=current_dir).check_returncode()
        subprocess.run(['nmake'], cwd=current_dir).check_returncode()
        exe = os.path.join(current_dir, "dist", args.name + ".exe")
        if args.win_sign:
            subprocess.run(f'{args.win_sign} {exe}', shell=True, cwd=current_dir)
        else:
            print(f'not win_sign')
        dest_dir = os.path.join(current_dir, "dist", args.name)
        os.makedirs(dest_dir, exist_ok=True)
        dest = os.path.join(dest_dir, args.name + ".exe")
        shutil.copy2(exe, dest)
        os.remove(exe)

        subprocess.run([f'{args.qt_bin}/windeployqt', dest,
                        "--dir=" + dest_dir,
                        "--release",
                        "--no-translations",
                        "--no-virtualkeyboard",
                        "--no-compiler-runtime"
                        ], cwd=current_dir).check_returncode()
        

        shutil.make_archive(os.path.join(current_dir, "dist", args.name), 'zip', dest_dir)
        with open(f"dist/{args.name}.zip.md5", 'w') as file:
            file.write(f"{args.name}.zip: " +  hashlib.md5(open(f"dist/{args.name}.zip", 'rb').read()).hexdigest())
          
    # 删除临时文件
    os.remove('webinterceptor_build.pro')
    if not isWin:
        os.remove('Info_build.plist')


if __name__ == '__main__':
    main()
