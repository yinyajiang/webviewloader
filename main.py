import webview
import threading
import time
import argparse
import json

def get_info(window, ua =None, wait_elements=None, wait_cookies=None):
    # wait for elements
    if wait_elements:   
        for name in wait_elements:
            el = window.dom.get_element(name)
            if not el:
                return None
        
    time.sleep(1)
    if wait_cookies:
        cookies = window.get_cookies()    
        for name in wait_cookies:
            if not any(name+"=" in c.output().strip() for c in cookies):
                return None

    cookies = window.get_cookies()

    nm = {}
    for c in cookies:
        k, v = c.output(header="").strip().split(";")[0].strip().split("=", 1)
        nm[k] = v
    info = {
        "ua": ua,
        "cookies": nm
    }
    print(json.dumps(info))
    return True



def hook(window, ua, element_names, cookie_names):
    def timer_func():
        info = get_info(window, ua, element_names, cookie_names)
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
    parser.add_argument('--title', default='',
                        help='Window title')
    parser.add_argument('--ua', default='Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
                        help='User agent string')
    parser.add_argument('--elements', nargs='+', default=[],
                        help='Element names to search for (can specify multiple)')
    parser.add_argument('--cookies', nargs='+', default=[],
                        help='Cookie names to search for (can specify multiple)')
    parser.add_argument('url', help='Target URL')
  
    args = parser.parse_args()

    window = webview.create_window(args.title, args.url)
    webview.start(lambda w: hook(w, ua=args.ua, element_names=args.elements, cookie_names=args.cookies), window, user_agent=args.ua)
