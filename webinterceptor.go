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

type WebInterceptorConfig struct {
	WinWebInterceptorAppURI    string
	WinWebInterceptorAppMd5URI string

	MacWebInterceptorAppURI    string
	MacWebInterceptorAppMd5URI string

	WebInterceptorAppWorkDir string
	WebInterceptorAppName    string

	CustomDownloadFileFunc func(url string, path string) error `json:"-"`
}

type WebInterceptorOptions struct {
	UA          string
	Title       string
	Width       int
	Height      int
	Banner      string
	BannerColor string
	ShowAddress bool
}

type WebInterceptorResult struct {
	UA      string            `json:"ua"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type WebInterceptor struct {
	cfg                WebInterceptorConfig
	lock               sync.Mutex
	webInterceptorPath string
}

func NewWebInterceptor(cfg WebInterceptorConfig) *WebInterceptor {
	if cfg.WebInterceptorAppName == "" {
		if isWindows() {
			cfg.WebInterceptorAppName = findBaseName(cfg.WinWebInterceptorAppURI)
		} else {
			cfg.WebInterceptorAppName = findBaseName(cfg.MacWebInterceptorAppURI)
		}
	}
	return &WebInterceptor{cfg: cfg}
}

func (l *WebInterceptor) HasMustCfg() bool {
	if isWindows() {
		return l.cfg.WinWebInterceptorAppURI != ""
	}
	return l.cfg.MacWebInterceptorAppURI != ""
}

func (l *WebInterceptor) CheckEnv(checkUpdate bool) (err error) {
	_, _, err = l.getWebInterceptorPath(checkUpdate)
	return
}

func (l *WebInterceptor) InstallEnv(checkUpdate bool, saveOpt WebInterceptorOptions) (err error) {
	err = l.installWebInterceptor(checkUpdate)
	if err != nil {
		return
	}
	l.saveOptions(saveOpt)
	return
}

func (l *WebInterceptor) Start(url string, opt WebInterceptorOptions) (result WebInterceptorResult, err error) {
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
	if opt.UA != "" {
		args = append(args, "--ua", opt.UA)
	}
	if opt.Banner != "" {
		args = append(args, "--banner", opt.Banner)
	}
	if opt.BannerColor != "" {
		args = append(args, "--banner-color", opt.BannerColor)
	}
	if opt.ShowAddress {
		args = append(args, "--address")
	}

	webInterceptorPath, err := l.GetWebInterceptorPath()
	if err != nil {
		return
	}
	fmt.Printf("webinterceptor cmd: %s %v\n", webInterceptorPath, args)
	c := exec.Command(webInterceptorPath, args...)
	stdout, err := c.Output()
	if len(stdout) > 0 {
		obj, e := findJsonObject(stdout)
		if e == nil {
			err = json.Unmarshal(obj, &result)
			if err == nil && result.URL == "" {
				err = fmt.Errorf("error: %s", string(stdout))
			}
			if err == nil {
				if result.Headers["User-Agent"] != "" {
					result.UA = result.Headers["User-Agent"]
				}
			}
			return
		}
	}
	return
}

func (l *WebInterceptor) GetWebInterceptorPath() (path string, err error) {
	path, _, err = l.getWebInterceptorPath(false)
	return
}

func (l *WebInterceptor) installWebInterceptor(checkUpdate bool) (err error) {
	releaser, err := l.getGlobalMutexLock()
	if err != nil {
		return
	}
	defer mutexRelease(releaser)
	_, _, err = l.getWebInterceptorPath(checkUpdate)
	return
}

func (l *WebInterceptor) getGlobalMutexLock() (releaser mutex.Releaser, err error) {
	releaser, err = mutexAcquire("install-"+l.cfg.WebInterceptorAppName, time.Minute*30)
	return releaser, err
}

func (l *WebInterceptor) getWebInterceptorPath(checkUpdate bool) (path string, useLast bool, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.webInterceptorPath != "" && fileutil.IsExist(l.webInterceptorPath) {
		return l.webInterceptorPath, true, nil
	}

	defer func() {
		if err == nil {
			if !isWindows() {
				path = filepath.Join(path, "Contents", "MacOS", l.cfg.WebInterceptorAppName)
			}
			l.webInterceptorPath = path
		}
	}()

	webInterceptorAppPath := filepath.Join(l.cfg.WebInterceptorAppWorkDir, l.cfg.WebInterceptorAppName)
	if isWindows() {
		webInterceptorAppPath = filepath.Join(webInterceptorAppPath, l.cfg.WebInterceptorAppName+".exe")
	} else {
		webInterceptorAppPath += ".app"
	}

	md5Url := l.cfg.WinWebInterceptorAppMd5URI
	if !isWindows() {
		md5Url = l.cfg.MacWebInterceptorAppMd5URI
	}

	netmd5 := ""
	loacalMd5Path := filepath.Join(l.cfg.WebInterceptorAppWorkDir, l.cfg.WebInterceptorAppName+".md5")
	exist := false
	if fileutil.IsExist(webInterceptorAppPath) {
		if md5Url == "" || !checkUpdate {
			return webInterceptorAppPath, true, nil
		}
		netmd5, err = downloadString(md5Url)
		if err != nil {
			fmt.Printf("download md5 failed: %v, %s\n", err, md5Url)
			return webInterceptorAppPath, true, nil
		}
		md5, _ := fileutil.ReadFileToString(loacalMd5Path)
		if strings.TrimSpace(md5) == strings.TrimSpace(netmd5) {
			return webInterceptorAppPath, true, nil
		}
		exist = true
	}

	url := l.cfg.MacWebInterceptorAppURI
	if isWindows() {
		url = l.cfg.WinWebInterceptorAppURI
	}

	tempPath := filepath.Join(l.cfg.WebInterceptorAppWorkDir, l.cfg.WebInterceptorAppName+".temp")
	os.Remove(tempPath)
	err = downloadFile(url, tempPath)
	if err != nil {
		fmt.Printf("download webview failed: %v, %s\n", err, url)
		if exist {
			err = nil
		}
		return webInterceptorAppPath, true, err
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
		os.RemoveAll(filepath.Dir(webInterceptorAppPath))
		unZip(tempPath, filepath.Dir(webInterceptorAppPath))
	} else {
		os.RemoveAll(webInterceptorAppPath)
		unZip(tempPath, filepath.Dir(webInterceptorAppPath))
		os.Chmod(webInterceptorAppPath, 0755)
	}
	os.Remove(tempPath)

	if !fileutil.IsExist(webInterceptorAppPath) {
		err = fmt.Errorf("webinterceptor app not found: %s", webInterceptorAppPath)
	}
	return webInterceptorAppPath, false, err
}

func (l *WebInterceptor) saveOptions(opt WebInterceptorOptions) (err error) {
	j, err := json.Marshal(opt)
	if err != nil {
		return
	}
	return os.WriteFile(filepath.Join(l.cfg.WebInterceptorAppWorkDir, l.cfg.WebInterceptorAppName+"_opt.json"), j, 0644)
}

func (l *WebInterceptor) loadAndMergeOptions(opt *WebInterceptorOptions) (err error) {
	if opt == nil {
		return
	}
	j, err := os.ReadFile(filepath.Join(l.cfg.WebInterceptorAppWorkDir, l.cfg.WebInterceptorAppName+"_opt.json"))
	if err != nil {
		return
	}
	var tmp WebInterceptorOptions
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
	if opt.Banner == "" {
		opt.Banner = tmp.Banner
	}
	if opt.BannerColor == "" {
		opt.BannerColor = tmp.BannerColor
	}
	if !opt.ShowAddress {
		opt.ShowAddress = tmp.ShowAddress
	}
	return
}
