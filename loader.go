package webviewloader

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/duke-git/lancet/v2/fileutil"
)

type Config struct {
	WindowsWebviewURI               string
	WindowsWebviewMd5URI            string
	WindowsDependniesComponentURI32 string
	WindowsDependniesComponentURI64 string

	MacWebviewURI    string
	MacWebviewMd5URI string

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
	Hidden       bool
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

func (l *Loader) HasMustCfg() bool {
	if isWindows() {
		return l.cfg.WindowsWebviewURI != ""
	}
	return l.cfg.MacWebviewURI != ""
}

func (l *Loader) CheckEnv() (err error) {
	err = checkComponent()
	return
}

func (l *Loader) InstallEnv() (err error) {
	err = checkComponent()
	if err != nil && isWindows() {
		err = installComponent(l.cfg.WindowsDependniesComponentURI32, l.cfg.WindowsDependniesComponentURI64, l.cfg.WebviewWorkDir)
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
	if opt.Hidden {
		args = append(args, "--hidden")
	}
	webviewPath, err := l.GetWebviewPath()
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

func (l *Loader) GetWebviewPath() (path string, err error) {
	path, _, err = l.getWebviewPath()
	return
}

func (l *Loader) getWebviewPath() (path string, useLast bool, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.webviewPath != "" {
		return l.webviewPath, true, nil
	}

	defer func() {
		if err == nil {
			if !isWindows() {
				path = filepath.Join(path, "Contents", "MacOS", l.cfg.WebviewName)
			}
			l.webviewPath = path
		}
	}()

	webviewPath := filepath.Join(l.cfg.WebviewWorkDir, l.cfg.WebviewName)
	if isWindows() {
		webviewPath += ".exe"
	} else {
		webviewPath += ".app"
	}

	md5Url := l.cfg.WindowsWebviewMd5URI
	if !isWindows() {
		md5Url = l.cfg.MacWebviewMd5URI
	}

	netmd5 := ""
	loacalMd5Path := filepath.Join(l.cfg.WebviewWorkDir, l.cfg.WebviewName+".md5")
	exist := false
	if fileutil.IsExist(webviewPath) {
		if md5Url == "" {
			return webviewPath, true, nil
		}
		netmd5, err = downloadString(md5Url)
		if err != nil {
			fmt.Printf("download md5 failed: %v, %s\n", err, md5Url)
			return webviewPath, true, nil
		}
		md5, _ := fileutil.ReadFileToString(loacalMd5Path)
		if strings.TrimSpace(md5) == strings.TrimSpace(netmd5) {
			return webviewPath, true, nil
		}
		exist = true
	}

	url := l.cfg.WindowsWebviewURI
	if !isWindows() {
		url = l.cfg.MacWebviewURI
	}

	tempPath := webviewPath + ".temp"
	err = downloadFile(url, tempPath)
	if err != nil {
		fmt.Printf("download webview failed: %v, %s\n", err, url)
		if exist {
			err = nil
		}
		return webviewPath, true, err
	}

	if !isWindows() {
		exec.Command("xattr", "-cr", tempPath).Run()
	}

	if netmd5 == "" {
		netmd5, _ = downloadString(md5Url)
	}
	if netmd5 != "" {
		os.WriteFile(loacalMd5Path, []byte(netmd5), 0644)
	}

	if isWindows() {
		os.Remove(webviewPath)
		err = os.Rename(tempPath, webviewPath)
	} else {
		os.RemoveAll(webviewPath)
		fileutil.UnZip(tempPath, filepath.Dir(webviewPath))
		os.Chmod(webviewPath, 0755)
	}

	return webviewPath, false, err
}
