import re
import subprocess
import argparse


def check_cert_valid(cert): 
    output = subprocess.check_output(f'security find-certificate -c "{cert}"', shell=True).decode('utf-8')
    if cert not in output:
        raise Exception(f'{cert} not found')
    return True

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Build Qt project')
    parser.add_argument('--cert',
                      help='maccert')
    parser.add_argument('--app',
                      help='app path')
    args = parser.parse_args()
    check_cert_valid(args.cert, args.app)
    subprocess.run([f'codesign', '--timestamp', '--force', '--deep', '--verify', '--verbose', '--sign', args.cert, f'{args.app}']).check_returncode()
    print(f'codesign success')