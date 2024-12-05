package webviewloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/elastic/go-windows"
	"github.com/wailsapp/go-webview2/webviewloader"

	"github.com/winlabs/gowin32"
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

func installComponent(cfg Config) (err error) {
	lock, err := newMutexLock(replaceMutexName("installwebview2"))
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
	}

	url := selectURI(selectURISt{
		x64:      cfg.WinDependniesComponentURI,
		x64lower: cfg.WinDependniesComponentLowerURI,
		x86:      cfg.WinDependniesComponentURIx86,
		x86lower: cfg.WinDependniesComponentLowerURIx86,
	})
	name := findName(url)
	dest := filepath.Join(cfg.WebviewAppWorkDir, "webview2", name)
	if cfg.CustomDownloadFileFunc != nil {
		err = cfg.CustomDownloadFileFunc(url, dest)
	} else {
		err = downloadFile(url, dest)
	}
	if err != nil {
		return err
	}
	err = exec.Command(dest, "/silent", "/install").Run()

	count := 0
	for count < 5 {
		e := checkComponent()
		if e == nil {
			os.Remove(dest)
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

func isWindow10OrGreater() bool {
	ok, err := gowin32.IsWindows10OrGreater()
	if err != nil {
		return false
	}
	return ok
}

func isBit32System() bool {
	info, err := windows.GetNativeSystemInfo()
	if err == nil {
		arch := info.ProcessorArchitecture
		if arch != 6 && arch != 9 && arch != 12 {
			return true
		}
	}
	return false
}

type selectURISt struct {
	x64, x64lower, x86, x86lower string
}

func selectURI(st selectURISt) string {
	if isBit32System() {
		if isWindow10OrGreater() {
			fmt.Println("use 32bit:", st.x86)
			return st.x86
		} else {
			fmt.Println("use lower 32bit: ", st.x86lower)
			return st.x86lower
		}
	} else {
		if isWindow10OrGreater() {
			fmt.Println("use 64bit:", st.x64)
			return st.x64
		} else {
			fmt.Println("use lower 64bit:", st.x64lower)
			return st.x64lower
		}
	}
}
