package webviewloader

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

const webInterceptorTestName = "TEST_WEBINTERCEPTOR"

func buildWebInterceptor(t *testing.T) {
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	os.Chdir(filepath.Join(oldDir, "webinterceptor"))
	var cmd *exec.Cmd
	if isWindows() {
		vsbat := "D:/vs2022/VC/Auxiliary/Build/vcvars64.bat"
		qtBin := "D:/Qt6.5.3/6.5.3/msvc2019_64/bin"
		cmd = exec.Command("python", "build.py", "--name", webInterceptorTestName, "--win-vsbat", vsbat, "--qt-bin", qtBin)
	} else {
		qtBin := "/Volumes/extern-usb/Apps/Qt/6.7.3/macos/bin"
		cmd = exec.Command("python", "build.py", "--name", webInterceptorTestName, "--qt-bin", qtBin)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	os.Chdir(oldDir)
}

func TestWebInterceptor(m *testing.T) {
	localDir, err := os.Getwd()
	if err != nil {
		m.Fatal(err)
	}

	cfg := WebInterceptorConfig{
		WinWebInterceptorAppURI:    filepath.Join(localDir, "webinterceptor/dist/"+webInterceptorTestName+".zip"),
		WinWebInterceptorAppMd5URI: filepath.Join(localDir, "webinterceptor/dist/"+webInterceptorTestName+".zip.md5"),

		MacWebInterceptorAppURI:    filepath.Join(localDir, "webinterceptor/dist/"+webInterceptorTestName+".app.zip"),
		MacWebInterceptorAppMd5URI: filepath.Join(localDir, "webinterceptor/dist/"+webInterceptorTestName+".app.zip.md5"),

		WebInterceptorAppWorkDir: filepath.Join(localDir, webInterceptorTestName),
		WebInterceptorAppName:    webInterceptorTestName,
	}
	os.RemoveAll(cfg.WebInterceptorAppWorkDir)

	l := NewWebInterceptor(cfg)

	build := false
	if isWindows() {
		m.Logf("win uri: %s", l.cfg.WinWebInterceptorAppURI)
		if !fileutil.IsExist(l.cfg.WinWebInterceptorAppURI) {
			build = true
		}
	} else {
		m.Logf("mac uri: %s", l.cfg.MacWebInterceptorAppURI)
		if !fileutil.IsExist(l.cfg.MacWebInterceptorAppURI) {
			build = true
		}
	}
	if build {
		buildWebInterceptor(m)
	}
	firstPath, useLast, err := l.getWebInterceptorPath(true)
	if err != nil {
		m.Fatal(err)
	}
	if useLast {
		m.Fatal("should not use last")
	}
	if err := l.InstallEnv(true); err != nil {
		m.Fatal(err)
	}

	info, err := l.Start("https://ww4.fmovies.co/film/gladiator-ii-1630857926/", WebInterceptorOptions{
		ShowAddress: true,
	})
	if err != nil {
		m.Fatal(err)
	}
	if info.URL == "" {
		m.Fatal("url is empty")
	}

	l2 := NewWebInterceptor(cfg)
	secondPath, useLast, err := l2.getWebInterceptorPath(true)
	if err != nil {
		m.Fatal(err)
	}
	if !useLast {
		m.Fatal("should use last")
	}
	if firstPath != secondPath {
		m.Fatal("path not equal")
	}

	releaser, err := l2.getGlobalMutexLock()
	if err != nil {
		m.Fatal(err)
	}
	releaser.Release()
}
