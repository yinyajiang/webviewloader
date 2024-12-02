package loadcookie

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
)

type Config struct {
	WindowsWebviewURL               string
	WindowsWebviewMd5URL            string
	WindowsDependniesComponentURL32 string
	WindowsDependniesComponentURL64 string

	MacWebviewURL    string
	MacWebviewMd5URL string

	WebviewWorkDir string
	WebviewName    string
}

type WebviewOptions struct {
	UA           string
	Title        string
	Width        int
	Height       int
	WaitElements []string
	WaitCookies  []string
}

type WebviewResult struct {
	UA      string            `json:"ua"`
	Cookies map[string]string `json:"cookies"`
}

type Loader struct {
	cfg         Config
	lock        sync.Mutex
	webviewPath string
}

func New(cfg Config) *Loader {
	if cfg.WebviewName == "" {
		cfg.WebviewName = "webview"
	}
	return &Loader{cfg: cfg}
}

func (l *Loader) CheckEnv() (err error) {
	err = checkComponent()
	return
}

func (l *Loader) InstallEnv() (err error) {
	err = checkComponent()
	if err != nil && isWindows() {
		err = installComponent(l.cfg.WindowsDependniesComponentURL32, l.cfg.WindowsDependniesComponentURL64, l.cfg.WebviewWorkDir)
	}
	return
}

func (l *Loader) Start(url string, opt WebviewOptions) (result WebviewResult, err error) {
	err = l.CheckEnv()
	if err != nil {
		return
	}

	args := []string{url}
	if opt.Title != "" {
		args = append(args, "--title", opt.Title)
	}
	if opt.Width > 0 {
		args = append(args, "--width", fmt.Sprintf("%d", opt.Width))
	}
	if opt.Height > 0 {
		args = append(args, "--height", fmt.Sprintf("%d", opt.Height))
	}
	if len(opt.WaitElements) > 0 {
		args = append(args, "--elements")
		args = append(args, opt.WaitElements...)
	}
	if len(opt.WaitCookies) > 0 {
		args = append(args, "--cookies")
		args = append(args, opt.WaitCookies...)
	}
	webviewPath, err := l.getWebviewPath()
	if err != nil {
		return
	}
	c := exec.Command(webviewPath, args...)
	stdout, err := c.Output()
	if err != nil {
		return
	}
	err = json.Unmarshal(stdout, &result)
	return
}

func (l *Loader) getWebviewPath() (path string, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.webviewPath != "" {
		return l.webviewPath, nil
	}

	defer func() {
		if err == nil {
			l.webviewPath = path
		}
	}()

	webviewPath := filepath.Join(l.cfg.WebviewWorkDir, l.cfg.WebviewName)
	if isWindows() {
		webviewPath += ".exe"
	}

	md5Url := l.cfg.WindowsWebviewMd5URL
	if !isWindows() {
		md5Url = l.cfg.MacWebviewMd5URL
	}

	netmd5 := ""
	loacalMd5Path := filepath.Join(l.cfg.WebviewWorkDir, l.cfg.WebviewName+".md5")
	exist := false
	if fileutil.IsExist(webviewPath) {
		if md5Url == "" {
			return webviewPath, nil
		}
		netmd5, err = downloadString(md5Url)
		if err != nil {
			fmt.Printf("download md5 failed: %v, %s\n", err, md5Url)
			return webviewPath, nil
		}
		md5, _ := fileutil.ReadFileToString(loacalMd5Path)
		if md5 == netmd5 {
			return webviewPath, nil
		}
		exist = true
	}

	url := l.cfg.WindowsWebviewURL
	if !isWindows() {
		url = l.cfg.MacWebviewURL
	}
	err = downloadFile(url, webviewPath+".temp")
	if err != nil {
		fmt.Printf("download webview failed: %v, %s\n", err, url)
		if exist {
			err = nil
		}
		return webviewPath, err
	}
	if netmd5 == "" {
		netmd5, _ = downloadString(md5Url)
	}
	if netmd5 != "" {
		os.WriteFile(loacalMd5Path, []byte(netmd5), 0644)
	}
	os.Remove(webviewPath)
	err = os.Rename(webviewPath+".temp", webviewPath)
	return webviewPath, err
}
