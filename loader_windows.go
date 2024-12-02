package loadcookie

import (
	"errors"
	"fmt"
	"github.com/elastic/go-windows"
	"github.com/wailsapp/go-webview2/webviewloader"
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
	url := url64
	info,err := windows.GetNativeSystemInfo()
	if err == nil {
		arch := info.ProcessorArchitecture
		if arch != 6 && arch != 9 && arch != 12 {
			url = url32
		}
	}
	downloadFile(url, "")
	
	return nil
}
