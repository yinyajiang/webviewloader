package webviewloader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWebInterceptor(m *testing.T) {
	localDir, err := os.Getwd()
	if err != nil {
		m.Fatal(err)
	}

	cfg := WebInterceptorConfig{
		WinWebInterceptorAppURI:    "http://10.2.51.27/dist/WebInterceptor.zip",
		WinWebInterceptorAppMd5URI: "http://10.2.51.27/dist/WebInterceptor.zip.md5",
		MacWebInterceptorAppURI:    "http://10.2.51.27/bin/WebInterceptor.app.zip",
		MacWebInterceptorAppMd5URI: "http://10.2.51.27/bin/WebInterceptor.app.zip.md5",
		WebInterceptorAppWorkDir:   filepath.Join(localDir, "Test"),
		WebInterceptorAppName:      "",
	}
	l := NewWebInterceptor(cfg)
	err = l.CheckEnv(false)
	if err != nil {
		m.Fatal(err)
	}
	l.webInterceptorPath = ""
	err = l.CheckEnv(true)
	if err != nil {
		m.Fatal(err)
	}
	l.webInterceptorPath = ""
	err = l.CheckEnv(true)
	if err != nil {
		m.Fatal(err)
	}

	info, err := l.Start("https://ww4.fmovies.co/film/gladiator-ii-1630857926/", WebInterceptorOptions{
		Title: "TEST",
	})
	if err != nil {
		m.Fatal(err)
	}
	m.Logf("info: %+v", info)
}
