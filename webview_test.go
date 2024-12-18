package webviewloader

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

const testName = "TEST_WEBVIEW"

func buildWebview(t *testing.T) {
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

func testFindjson(m *testing.T) {
	obj, e := findJsonObject([]byte(`
testestetsetetsgadsg
{
    "_id": "675a941634445a1baf73aff4",
    "index": 0,
    "guid": "e45351cc-1f0f-49e7-9f53-8c730da216a4",
    "isActive": true,
    "balance": "$3,813.09",
    "picture": "http://placehold.it/32x32",
    "age": 20,
    "eyeColor": "brown",
    "name": "Monroe Church",
    "gender": "male",
    "company": "DIGIFAD",
    "email": "monroechurch@digifad.com",
    "phone": "+1 (857) 567-3283",
    "address": "504 Thatford Avenue, Newcastle, West Virginia, 8256",
    "about": "Adipisicing minim ullamco culpa exercitation dolor. Fugiat nisi proident in minim proident qui enim ullamco voluptate qui mollit. Ut incididunt laborum duis est ipsum ad ex voluptate non. Proident laboris quis dolore qui elit consequat est cupidatat aute veniam aliquip.\r\n",
    "registered": "2015-04-07T12:01:11 -08:00",
    "latitude": 77.047377,
    "longitude": -101.104828,
    "tags": [
      "commodo",
      "esse",
      "commodo",
      "amet",
      "reprehenderit",
      "cillum",
      "laborum"
    ],
    "friends": [
      {
        "id": 0,
        "name": "Peggy Ewing"
      },
      {
        "id": 1,
        "name": "Alana Gonzalez"
      },
      {
        "id": 2,
        "name": "Blanca Joyce"
      }
    ],
    "greeting": "Hello, Monroe Church! You have 7 unread messages.",
    "favoriteFruit": "strawberry"
  }
testestetsetetsgadsgsgasdg
	`))
	if e != nil {
		m.Fatal(e)
	}
	var result map[string]interface{}
	e = json.Unmarshal(obj, &result)
	if e != nil {
		m.Fatal(e)
	}
	m.Logf("obj: %v", result)
	if len(result) == 0 {
		m.Fatal("obj is empty")
	}

}

func TestMain(m *testing.T) {
	testFindjson(m)

	localDir, err := os.Getwd()
	if err != nil {
		m.Fatal(err)
	}
	os.RemoveAll(filepath.Join(localDir, "Test"))

	cfg := WebviewConfig{
		WinWebviewAppURI:    filepath.Join(localDir, "webview/dist/"+testName+".exe"),
		WinWebviewAppMd5URI: filepath.Join(localDir, "webview/dist/"+testName+".exe.md5"),

		MacWebviewAppURI:    filepath.Join(localDir, "webview/dist/"+testName+".app.zip"),
		MacWebviewAppMd5URI: filepath.Join(localDir, "webview/dist/"+testName+".app.zip.md5"),

		WebviewAppWorkDir: filepath.Join(localDir, "Test"),
		WebviewAppName:    testName,
	}
	l := NewWebview(cfg)

	build := false
	if isWindows() {
		m.Logf("win uri: %s", l.cfg.WinWebviewAppURI)
		if !fileutil.IsExist(l.cfg.WinWebviewAppURI) {
			build = true
		}
	} else {
		m.Logf("mac uri: %s", l.cfg.MacWebviewAppURI)
		if !fileutil.IsExist(l.cfg.MacWebviewAppURI) {
			build = true
		}
	}
	if build {
		buildWebview(m)
	}
	firstPath, useLast, err := l.getWebviewPath(true)
	if err != nil {
		m.Fatal(err)
	}
	if useLast {
		m.Fatal("should not use last")
	}
	if err := l.InstallEnv(true); err != nil {
		m.Fatal(err)
	}

	info, err := l.Start("https://www.baidu.com", WebviewOptions{})
	if err != nil {
		m.Fatal(err)
	}
	if info.UA == "" || len(info.Cookies) == 0 {
		m.Fatal("ua or cookies is empty")
	}

	l2 := NewWebview(cfg)
	secondPath, useLast, err := l2.getWebviewPath(true)
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
