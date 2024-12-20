import re
import subprocess
import argparse

def get_cert(): 
    prename='Developer ID Application:'
    output = subprocess.check_output(f'security find-certificate -c "{prename}"', shell=True).decode('utf-8')
    match = re.compile(f'"({prename}.+)"').search(output)
    cert = match.group(1)
    if cert == '':
        raise Exception('No certificate found')
    return cert

def check_cert_valid(cert): 
    prename='Developer ID Application:'
    output = subprocess.check_output(f'security find-certificate -c "{prename}"', shell=True).decode('utf-8')
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