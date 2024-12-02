package loadcookie

import (
	"errors"
	"fmt"
	"github.com/elastic/go-windows"
	"github.com/wailsapp/go-webview2/webviewloader"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func checkComponent() (err error) {
	ver, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString("")
	if err != nil || ver == "" {
		if err == nil {
			err = errors.New("webview2 not found")
		}
		return err
	}
	fmt.Println("webview2 version:", ver)
	return nil
}

func installComponent(url32, url64, cacheDir string) (err error) {
	if url32 ==""{
		url32= "https://github.com/yinyajiang/load_cookie/releases/download/webview2/MicrosoftEdgeWebView2RuntimeInstallerX86.exe"
	}
	if url64 == "" {
		url64 = "https://github.com/yinyajiang/load_cookie/releases/download/webview2/MicrosoftEdgeWebView2RuntimeInstallerX64.exe"
	}

	url := url64
	info,err := windows.GetNativeSystemInfo()
	if err == nil {
		arch := info.ProcessorArchitecture
		if arch != 6 && arch != 9 && arch != 12 {
			url = url32
		}
	}

	arr := strings.Split(url, "/")
	name := arr[len(arr)-1]
	dest := filepath.Join(cacheDir, "webview2", name)
	err = downloadFile(url, dest)
	if err != nil {
		return err
	}
	err = exec.Command(dest, "/silent", "/install").Run()

	count := 0
	for count < 5 {
		e := checkComponent()
		if e == nil {
			return nil
		}
		time.Sleep(time.Second * 1)
		count++
	}
	return err
}
