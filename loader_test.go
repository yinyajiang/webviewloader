package webviewloader

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

const testName = "TEST_WEBVIEW"

func TestBuildWebview(t *testing.T) {
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	os.Chdir(filepath.Join(oldDir, "webview"))
	var cmd *exec.Cmd
	if isWindows() {
		cmd = exec.Command("cmd", "/C", "build.bat", "--name", testName)

	} else {
		cmd = exec.Command("bash", "build.sh", "--name", testName)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	os.Chdir(oldDir)
}

func TestMain(m *testing.T) {
	localDir, err := os.Getwd()
	if err != nil {
		m.Fatal(err)
	}
	os.RemoveAll(filepath.Join(localDir, "Test"))

	cfg := Config{
		WindowsWebviewURI:               filepath.Join(localDir, "webview/dist/"+testName+".exe"),
		WindowsWebviewMd5URI:            filepath.Join(localDir, "webview/dist/"+testName+".exe.md5"),
		WindowsDependniesComponentURI32: "",
		WindowsDependniesComponentURI64: "",

		MacWebviewURI:    filepath.Join(localDir, "webview/dist/"+testName+".app.zip"),
		MacWebviewMd5URI: filepath.Join(localDir, "webview/dist/"+testName+".app.zip.md5"),

		WebviewWorkDir: filepath.Join(localDir, "Test"),
		WebviewName:    testName,
	}
	l := New(cfg)

	build := false
	if isWindows() {
		if !fileutil.IsExist(l.cfg.WindowsWebviewURI) {
			build = true
		}
	} else {
		if !fileutil.IsExist(l.cfg.MacWebviewURI) {
			build = true
		}
	}
	if build {
		TestBuildWebview(m)
	}

	if err := l.InstallEnv(); err != nil {
		m.Fatal(err)
	}
	firstPath, useLast, err := l.getWebviewPath()
	if err != nil {
		m.Fatal(err)
	}
	if useLast {
		m.Fatal("should not use last")
	}

	info, err := l.Start("https://www.baidu.com", WebviewOptions{})
	if err != nil {
		m.Fatal(err)
	}
	if info.UA == "" || len(info.Cookies) == 0 {
		m.Fatal("ua or cookies is empty")
	}

	l2 := New(cfg)
	secondPath, useLast, err := l2.getWebviewPath()
	if err != nil {
		m.Fatal(err)
	}
	if !useLast {
		m.Fatal("should use last")
	}
	if firstPath != secondPath {
		m.Fatal("path not equal")
	}
}
