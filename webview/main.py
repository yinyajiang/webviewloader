import webview
import threading
import time
import argparse
import json
import sys
import os
import urllib.parse 
import datetime


class options:
    def __init__(self):
        self.ua = None
        self.wait_elements = None
        self.wait_cookies = None
        self.wait_domains = None
        self.write_cookies = False
        self.forever = False

def get_cookies_nm(window):
    cookies = window.get_cookies()
    nm = {}
    for c in cookies:
        k, v = c.output(header="").strip().split(";")[0].strip().split("=", 1)
        nm[k] = v
    return nm


def parse_cookie_expires(cookie):
    if cookie and 'expires' in cookie and cookie['expires']:
        if isinstance(cookie['expires'], int):
            return cookie['expires']
        try:
            return int(datetime.datetime.strptime(str(cookie['expires']), "%a, %d %b %Y %H:%M:%S %Z").timestamp())
        except Exception:
            pass
        try:
            return int(datetime.datetime.strptime(str(cookie['expires']), "%Y-%m-%d %H:%M:%S %z").timestamp())
        except Exception:
            pass
    return 2147483647
        

def write_cookies(window, filename):
    with open(filename, "w") as f:
        f.write("# Netscape HTTP Cookie File\n")
        f.write("# Created by webview\n\n")
        cookies = window.get_cookies()
        for c in cookies:
            for k in c.keys():
                cookie = c[k]
                name =cookie.key
                value = cookie.value
                path = cookie['path']
                domain = cookie['domain']
                expiry = parse_cookie_expires(cookie)
                secure = str(cookie['secure']).upper()
                include_subdomains = str(domain[0] == '.').upper()
                f.write(f"{domain}\t{include_subdomains}\t{path}\t{secure}\t{expiry}\t{name}\t{value}\n")
    return filename 


def output_info(window, opts=None):
    info = {
        "ua": opts.ua if opts else None,
        "cookies": get_cookies_nm(window),
        "url": window.get_current_url()
    }

    if opts.write_cookies:
        filename = write_cookies(window, opts.write_cookies)
        info["cookies_file"] = filename

    print(json.dumps(info)+"\n", flush=True)
    return True


def is_ready(window, opts=None):
    if not opts:
        return True
    
    # wait for elements
    if opts.wait_elements:   
        for name in opts.wait_elements:
            try:
                el = window.dom.get_element(name)
                if not el:
                    return False
            except Exception:
                return False
        
    # wait for cookies
    if opts.wait_cookies:
        cookies = window.get_cookies()
        for name in opts.wait_cookies:
            if not any(name+"=" in c.output().strip() for c in cookies):
                return False
            
    # wait for domains
    if opts.wait_domains:
        current_url = window.get_current_url()
        if not current_url:
            return False
        parsed_url = urllib.parse.urlparse(current_url)
        return any(d == parsed_url.netloc for d in opts.wait_domains)
    
    return True



def hook(window, opts):
    def timer_func():
        if is_ready(window, opts):
            outed = output_info(window, opts)
            if outed and not opts.forever:
                window.destroy()
                return
        start_timer()
        
    def start_timer():
        timer = threading.Timer(1, timer_func)
        timer.start()
    start_timer()


if __name__ == '__main__':
    isWin = sys.platform.lower().startswith('win')

    parser = argparse.ArgumentParser(description='Cookie loader with customizable parameters')
    parser.add_argument('url', help='Target URL')
    parser.add_argument('--title', default='',
                        help='Window title')
    parser.add_argument('--ua', default='Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36' if 
                        isWin else 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Safari/605.1.15',
                        help='User agent string')
    parser.add_argument('--elements', nargs='+', default=[],
                        help='Element names to search for (can specify multiple)')
    parser.add_argument('--cookies', nargs='+', default=[],
                        help='Cookie names to search for (can specify multiple)')
    parser.add_argument('--domains', nargs='+', default=[],
                        help='Domains to search for (can specify multiple)')
    parser.add_argument('--width', default=800, type=int,
                        help='Window width')
    parser.add_argument('--height', default=600, type=int,
                        help='Window height')
    parser.add_argument('--hidden', action='store_true',
                        help='Hide window')
    parser.add_argument('--forever', action='store_true',
                        help='Run forever')
    parser.add_argument('--write-cookies',
                        help='Write cookies to a Netscape HTTP Cookie File')
    args = parser.parse_args()

    title = args.title
    if not title and len(sys.argv) > 1:
        title = os.path.basename(sys.argv[0]).split(".")[0]

    opts = options()
    opts.ua = args.ua
    opts.wait_elements = args.elements
    opts.wait_cookies = args.cookies
    opts.wait_domains = args.domains
    opts.write_cookies = args.write_cookies
    opts.forever = args.forever
    window = webview.create_window(title, args.url, width=args.width, height=args.height, hidden=args.hidden)
    webview.start(lambda w: hook(w, opts), window, user_agent=args.ua)
