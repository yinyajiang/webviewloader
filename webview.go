package webviewloader

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/juju/mutex"
)

type WebviewConfig struct {
	WinWebviewAppURI       string
	WinWebviewAppMd5URI    string
	WinWebviewAppURIx86    string
	WinWebviewAppMd5URIx86 string

	WinWebviewAppLowerURI       string
	WinWebviewAppLowerMd5URI    string
	WinWebviewAppLowerURIx86    string
	WinWebviewAppLowerMd5URIx86 string

	WinDependniesComponentURI         string
	WinDependniesComponentURIx86      string
	WinDependniesComponentLowerURI    string
	WinDependniesComponentLowerURIx86 string

	MacWebviewAppURI    string
	MacWebviewAppMd5URI string

	WebviewAppWorkDir string
	WebviewAppName    string

	CustomDownloadFileFunc func(url string, path string) error `json:"-"`
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

type WebView struct {
	cfg         WebviewConfig
	lock        sync.Mutex
	webviewPath string
}

func NewWebview(cfg WebviewConfig) *WebView {
	if cfg.WebviewAppName == "" {
		if isWindows() {
			cfg.WebviewAppName = findBaseName(cfg.WinWebviewAppURI)
		} else {
			cfg.WebviewAppName = findBaseName(cfg.MacWebviewAppURI)
		}
	}
	if cfg.WinDependniesComponentLowerURIx86 == "" {
		cfg.WinDependniesComponentLowerURIx86 = "https://github.com/yinyajiang/webviewloader/releases/download/webview2/MicrosoftEdgeWebView2RuntimeInstallerLowX86.exe"
	}
	if cfg.WinDependniesComponentURIx86 == "" {
		cfg.WinDependniesComponentURIx86 = "https://github.com/yinyajiang/webviewloader/releases/download/webview2/MicrosoftEdgeWebView2RuntimeInstallerX86.exe"
	}
	if cfg.WinDependniesComponentLowerURI == "" {
		cfg.WinDependniesComponentLowerURI = "https://github.com/yinyajiang/webviewloader/releases/download/webview2/MicrosoftEdgeWebView2RuntimeInstallerLowX64.exe"
	}
	if cfg.WinDependniesComponentURI == "" {
		cfg.WinDependniesComponentURI = "https://github.com/yinyajiang/webviewloader/releases/download/webview2/MicrosoftEdgeWebView2RuntimeInstallerX64.exe"
	}
	if cfg.WinWebviewAppLowerURI == "" {
		cfg.WinWebviewAppLowerURI = cfg.WinWebviewAppURI
	}
	if cfg.WinWebviewAppLowerMd5URI == "" {
		cfg.WinWebviewAppLowerMd5URI = cfg.WinWebviewAppMd5URI
	}
	if cfg.WinWebviewAppLowerURIx86 == "" {
		cfg.WinWebviewAppLowerURIx86 = cfg.WinWebviewAppURIx86
	}
	if cfg.WinWebviewAppLowerMd5URIx86 == "" {
		cfg.WinWebviewAppLowerMd5URIx86 = cfg.WinWebviewAppMd5URIx86
	}
	return &WebView{cfg: cfg}
}

func (l *WebView) HasMustCfg() bool {
	if isWindows() {
		return l.cfg.WinWebviewAppURI != ""
	}
	return l.cfg.MacWebviewAppURI != ""
}

func (l *WebView) CheckEnv(checkUpdate bool) (err error) {
	err = checkWebviewComponent()
	if err == nil {
		_, _, err = l.getWebviewPath(checkUpdate)
	}
	return
}

func (l *WebView) InstallEnv(checkUpdate bool, opt WebviewOptions) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var componentErr error
	var webviewErr error
	go func() {
		defer wg.Done()
		componentErr = l.installWebviewComponent()
	}()
	go func() {
		defer wg.Done()
		webviewErr = l.installWebview(checkUpdate)
	}()
	wg.Wait()
	if componentErr != nil {
		return componentErr
	}
	if webviewErr != nil {
		return webviewErr
	}
	l.saveOptions(opt)
	return nil
}

func (l *WebView) Start(url string, opt WebviewOptions) (result WebviewResult, err error) {
	err = l.CheckEnv(false)
	if err != nil {
		return
	}
	l.loadAndMergeOptions(&opt)

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
	if opt.UA != "" {
		args = append(args, "--ua", opt.UA)
	}

	webviewPath, err := l.GetWebviewPath()
	if err != nil {
		return
	}
	fmt.Printf("webview cmd: %s %v\n", webviewPath, args)
	c := exec.Command(webviewPath, args...)
	stdout, err := c.Output()
	if len(stdout) > 0 {
		obj, e := findJsonObject(stdout)
		if e == nil {
			err = json.Unmarshal(obj, &result)
			return
		}
	}
	return
}

func (l *WebView) GetWebviewPath() (path string, err error) {
	path, _, err = l.getWebviewPath(false)
	return
}

func (l *WebView) installWebviewComponent() (err error) {
	err = checkWebviewComponent()
	if err != nil && isWindows() {
		err = installWebviewComponent(l.cfg)
	}
	return
}

func (l *WebView) installWebview(checkUpdate bool) (err error) {
	releaser, err := l.getGlobalMutexLock()
	if err != nil {
		return
	}
	defer mutexRelease(releaser)
	_, _, err = l.getWebviewPath(checkUpdate)
	return
}

func (l *WebView) getGlobalMutexLock() (releaser mutex.Releaser, err error) {
	releaser, err = mutexAcquire("install-"+l.cfg.WebviewAppName, time.Minute*10)
	return releaser, err
}

func (l *WebView) getWebviewPath(checkUpdate bool) (path string, useLast bool, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.webviewPath != "" && fileutil.IsExist(l.webviewPath) {
		return l.webviewPath, true, nil
	}

	defer func() {
		if err == nil {
			if !isWindows() {
				path = filepath.Join(path, "Contents", "MacOS", l.cfg.WebviewAppName)
			}
			l.webviewPath = path
		}
	}()

	webviewPath := filepath.Join(l.cfg.WebviewAppWorkDir, l.cfg.WebviewAppName)
	if isWindows() {
		webviewPath += ".exe"
	} else {
		webviewPath += ".app"
	}

	md5Url := l.cfg.MacWebviewAppMd5URI
	if isWindows() {
		md5Url = selectURI(selectURISt{
			x64:      l.cfg.WinWebviewAppMd5URI,
			x64lower: l.cfg.WinWebviewAppLowerMd5URI,
			x86:      l.cfg.WinWebviewAppMd5URIx86,
			x86lower: l.cfg.WinWebviewAppLowerMd5URIx86,
		})
	}

	netmd5 := ""
	loacalMd5Path := filepath.Join(l.cfg.WebviewAppWorkDir, l.cfg.WebviewAppName+".md5")
	exist := false
	if fileutil.IsExist(webviewPath) {
		if md5Url == "" || !checkUpdate {
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

	url := l.cfg.MacWebviewAppURI
	if isWindows() {
		url = selectURI(selectURISt{
			x64:      l.cfg.WinWebviewAppURI,
			x64lower: l.cfg.WinWebviewAppLowerURI,
			x86:      l.cfg.WinWebviewAppURIx86,
			x86lower: l.cfg.WinWebviewAppLowerURIx86,
		})
	}

	tempPath := webviewPath + ".temp"
	os.Remove(tempPath)
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
		unZip(tempPath, filepath.Dir(webviewPath))
		os.Chmod(webviewPath, 0755)
		os.Remove(tempPath)
	}

	return webviewPath, false, err
}

func (l *WebView) saveOptions(opt WebviewOptions) (err error) {
	j, err := json.Marshal(opt)
	if err != nil {
		return
	}
	return os.WriteFile(filepath.Join(l.cfg.WebviewAppWorkDir, l.cfg.WebviewAppName+"_opt.json"), j, 0644)
}

func (l *WebView) loadAndMergeOptions(opt *WebviewOptions) (err error) {
	if opt == nil {
		return
	}
	j, err := os.ReadFile(filepath.Join(l.cfg.WebviewAppWorkDir, l.cfg.WebviewAppName+"_opt.json"))
	if err != nil {
		return
	}
	var tmp WebviewOptions
	err = json.Unmarshal(j, &tmp)
	if err != nil {
		return
	}
	if opt.UA == "" {
		opt.UA = tmp.UA
	}
	if opt.Title == "" {
		opt.Title = tmp.Title
	}
	if opt.Width == 0 {
		opt.Width = tmp.Width
	}
	if opt.Height == 0 {
		opt.Height = tmp.Height
	}
	return
}
