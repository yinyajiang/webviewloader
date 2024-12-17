from PyQt6.QtWidgets import QApplication, QMainWindow, QLabel, QVBoxLayout, QWidget
from PyQt6.QtWebEngineCore import QWebEnginePage, QWebEngineProfile, QWebEngineUrlRequestInterceptor, QWebEngineSettings
from PyQt6.QtWebEngineWidgets import QWebEngineView
from PyQt6.QtCore import QUrl, Qt
import sys
import json
import argparse
import os


class UrlRequestInterceptor(QWebEngineUrlRequestInterceptor):
    def __init__(self):
        super().__init__()

    def interceptRequest(self, info):
        url = info.requestUrl().toString()
        if self._is_playable(url):
            headers = {}
            try:
                http_headers = info.httpHeaders()
                for header in http_headers.keys():
                    header_key = bytes(header).decode()
                    header_value = bytes(http_headers[header]).decode()
                    headers[header_key] = header_value
            except Exception as e:
                print(f"Error getting headers: {e}", file=sys.stderr, flush=True)
            obj = {
                "url": url,
                "method": bytes(info.requestMethod()).decode(),
                "headers": headers
            }
            print(json.dumps(obj)+"\n", flush=True)
            sys.exit(0)

    def _is_playable(self, url):
        return '.m3u8' in url.lower()


class WebInterceptor(QWebEngineProfile):
    def __init__(self, ua):
        super().__init__()
        self.interceptor = UrlRequestInterceptor()
        self.setUrlRequestInterceptor(self.interceptor)
        if ua:
            self.setHttpUserAgent(ua)
        # 修正：使用 QWebEngineSettings 的 PlaybackRequiresUserGesture
        self.settings().setAttribute(QWebEngineSettings.WebAttribute.PlaybackRequiresUserGesture, True)


class Browser(QMainWindow):
    def __init__(self, url, title, ua, width, height, banner, banner_color):
        super().__init__()
        self.setWindowTitle(title)
        
        # 创建中央窗口部件
        central_widget = QWidget()
        self.setCentralWidget(central_widget)
        
        # 创建垂直布局
        layout = QVBoxLayout(central_widget)
        
        self.banner = QLabel(banner)
        self.banner.setAlignment(Qt.AlignmentFlag.AlignCenter)
        self.banner.setStyleSheet("""
            QLabel {
                background-color: %s;
                color: white;
                padding: 6px 20px;
                font-size: 13px;
                font-weight: bold;
                border-radius: 3px;
                margin: 8px 15px;
                max-height: 22px;
                box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
                letter-spacing: 1px;
            }
        """ % banner_color)
        self.banner.setMaximumHeight(36)  # 增加最大高度
        
        # 创建web视图
        self.web_view = QWebEngineView()
        self.profile = WebInterceptor(ua)
        self.page = QWebEnginePage(self.profile, self)
        self.web_view.setPage(self.page)
        
        layout.addWidget(self.banner)
        layout.addWidget(self.web_view)
        
        # 设置布局的边距
        layout.setContentsMargins(0, 0, 0, 0)
        layout.setSpacing(0)
        
        self.resize(width, height)
        self.web_view.load(QUrl(url))

    def closeEvent(self, event):
        # 在窗口关闭时确保正确清理
        self.web_view.setPage(None)  # 解除页面引用
        if self.page:
            self.page.deleteLater()  # 删除页面
        super().closeEvent(event)
        obj = {
            "Error": "Closed by user.",
        }
        print(json.dumps(obj)+"\n", flush=True)
  


if __name__ == '__main__':
    isWin = sys.platform.lower().startswith('win')

    parser = argparse.ArgumentParser(description='Playable url interceptor')
    parser.add_argument('url', help='Target URL')
    parser.add_argument('--title', default='',
                        help='Window title')
    parser.add_argument('--banner', default='Please start playing the video first, and then it will start parsing.',
                        help='Banner text')
    parser.add_argument('--banner_color', default='#FF4E50',
                        help='Banner background color')
    parser.add_argument('--ua', default='')
    parser.add_argument('--width', default=1024, type=int,
                        help='Window width')
    parser.add_argument('--height', default=768, type=int,
                        help='Window height')
    args = parser.parse_args()

    if not args.title and len(sys.argv) > 1:
        args.title = os.path.basename(sys.argv[0]).split(".")[0]

    app = QApplication(sys.argv)
    browser = Browser(url=args.url,
                       ua=args.ua,
                       title=args.title, 
                       width=args.width,
                       height=args.height,
                       banner=args.banner,
                       banner_color=args.banner_color)
    browser.show()
    sys.exit(app.exec())
