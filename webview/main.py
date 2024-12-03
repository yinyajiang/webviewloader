import webview
import threading
import time
import argparse
import json
import sys
import os

def get_cookies_nm(window):
    cookies = window.get_cookies()
    nm = {}
    for c in cookies:
        k, v = c.output(header="").strip().split(";")[0].strip().split("=", 1)
        nm[k] = v
    return nm

def get_info(window, ua=None):
    info = {
        "ua": ua,
        "cookies": get_cookies_nm(window)
    }
    print(json.dumps(info)+"\n", flush=True)
    return True


def is_ready(window, wait_elements=None, wait_cookies=None):
    # wait for elements
    if wait_elements:   
        for name in wait_elements:
            el = window.dom.get_element(name)
            if not el:
                return False
        
    time.sleep(1)
    if wait_cookies:
        cookies = window.get_cookies()
        for name in wait_cookies:
            if not any(name+"=" in c.output().strip() for c in cookies):
                return False
    return True


def hook(window, ua, wait_elements, wait_cookies):
    def timer_func():
        if is_ready(window, wait_elements, wait_cookies):
            info = get_info(window, ua)
            if info:
                window.destroy()
                return
        start_timer()
        
    def start_timer():
        timer = threading.Timer(1, timer_func)
        timer.start()
    start_timer()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Cookie loader with customizable parameters')
    parser.add_argument('url', help='Target URL')
    parser.add_argument('--title', default='',
                        help='Window title')
    parser.add_argument('--ua', default='Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
                        help='User agent string')
    parser.add_argument('--elements', nargs='+', default=[],
                        help='Element names to search for (can specify multiple)')
    parser.add_argument('--cookies', nargs='+', default=[],
                        help='Cookie names to search for (can specify multiple)')
    parser.add_argument('--width', default=800, type=int,
                        help='Window width')
    parser.add_argument('--height', default=600, type=int,
                        help='Window height')
    parser.add_argument('--hidden', action='store_true',
                        help='Hide window')
    args = parser.parse_args()

    title = args.title
    if not title and len(sys.argv) > 1:
        title = os.path.basename(sys.argv[0]).split(".")[0]

    window = webview.create_window(title, args.url, width=args.width, height=args.height, hidden=args.hidden)
    webview.start(lambda w: hook(w, ua=args.ua, wait_elements=args.elements, wait_cookies=args.cookies), window, user_agent=args.ua)
