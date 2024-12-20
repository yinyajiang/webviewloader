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

def get_cert(): 
    prename='Developer ID Application:'
    output = subprocess.check_output(f'security find-certificate -c "{prename}"', shell=True).decode('utf-8')
    match = re.compile(f'"({prename}.+)"').search(output)
    cert = match.group(1)
    if cert == '':
        raise Exception('No certificate found')
    return cert


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
    parser.add_argument('--win-sign',
                      help='Path to win sign file')
    parser.add_argument('--win-vsbat',
                      help='Path to win vsbat file')
    parser.add_argument('--qt-bin',
                      help='Path to qt bin file')
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
                content = content.replace('{bundle_id}', args.bundle_id)
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
        subprocess.run([f'{args.qt_bin}/bin/macdeployqt', f'dist/{args.name}.app',
                        ], cwd=current_dir).check_returncode()
        cert = get_cert()
        if cert:
            subprocess.run([f'codesign', '--timestamp', '--force', '--deep', '--verify', '--verbose', '--sign', cert, f'dist/{args.name}.app'], cwd=current_dir).check_returncode()
            print(f'codesign success')
            subprocess.run(['zip', '-ry', f'{args.name}.app.zip', f'{args.name}.app'], cwd=os.path.join(current_dir, 'dist')).check_returncode()
            print(f'zip success')
            with open(f"dist/{args.name}.app.zip.md5", 'w') as file:
                file.write(f"{args.name}.app.zip: " +  hashlib.md5(open(f"dist/{args.name}.app.zip", 'rb').read()).hexdigest())
        else:
            print(f'not found cert')
    else:
        subprocess.run([f'{args.qt_bin}/qmake', 'webinterceptor_build.pro'], cwd=current_dir).check_returncode()
        subprocess.run(['nmake'], cwd=current_dir).check_returncode()
        exe = os.path.join(current_dir, "dist", args.name + ".exe")
        subprocess.run(f'{args.win_sign} {exe}', shell=True, cwd=current_dir)

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
