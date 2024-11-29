import re
import subprocess

def get_cert(): 
    prename='Developer ID Application:'
    output = subprocess.check_output(f'security find-certificate -c "{prename}"', shell=True).decode('utf-8')
    match = re.compile(f'"({prename}.+)"').search(output)
    cert = match.group(1)
    if cert == '':
        raise Exception('No certificate found')
    return cert

if __name__ == '__main__':
    print(get_cert())
