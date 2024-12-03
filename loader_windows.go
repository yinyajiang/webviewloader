package webviewloader

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/elastic/go-windows"
	"github.com/wailsapp/go-webview2/webviewloader"

	xwindows "golang.org/x/sys/windows"
)

func checkComponent() (err error) {
	ver, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString("")
	if err != nil || ver == "" {
		err = fmt.Errorf("webview2 not found, %v", err)
		return err
	}
	fmt.Println("webview2 version:", ver)
	return nil
}

func installComponent(url32, url64, cacheDir string) (err error) {
	lock, err := newMutexLock("__install_webview2")
	if err == nil {
		hasSleep := false
		for {
			err = lock.Lock()
			if err == nil {
				defer lock.Unlock()
				break
			}
			time.Sleep(time.Second * 1)
			hasSleep = true
		}
		if hasSleep {
			time.Sleep(time.Second * 2)
			//另一个进程可能安装完成
			err = checkComponent()
			if err == nil {
				return nil
			}
		}
		return nil
	}

	if url32 == "" {
		url32 = "https://github.com/yinyajiang/load_cookie/releases/download/webview2/MicrosoftEdgeWebView2RuntimeInstallerX86.exe"
	}
	if url64 == "" {
		url64 = "https://github.com/yinyajiang/load_cookie/releases/download/webview2/MicrosoftEdgeWebView2RuntimeInstallerX64.exe"
	}

	url := url64
	info, err := windows.GetNativeSystemInfo()
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

type mutexLock struct {
	name  string
	mutex xwindows.Handle
}

func newMutexLock(name string) (*mutexLock, error) {
	if !strings.HasPrefix(name, "Global\\") {
		name = "Global\\" + name
	}
	m, err := xwindows.CreateMutex(nil, false, xwindows.StringToUTF16Ptr(name))
	if err != nil {
		return nil, err
	}
	return &mutexLock{mutex: m, name: name}, nil
}

func (m *mutexLock) Lock() error {
	event, err := xwindows.WaitForSingleObject(m.mutex, xwindows.INFINITE)
	if err != nil {
		return err
	}
	if event == xwindows.WAIT_OBJECT_0 {
		return nil
	}
	return fmt.Errorf("mutex failed to lock: %s", m.name)
}

func (m *mutexLock) TryLock() error {
	event, err := xwindows.WaitForSingleObject(m.mutex, 0)
	if err != nil {
		return err
	}
	if event == xwindows.WAIT_OBJECT_0 {
		return nil
	}
	return fmt.Errorf("another has locked the mutex")
}

func (m *mutexLock) Unlock() error {
	return xwindows.ReleaseMutex(m.mutex)
}
