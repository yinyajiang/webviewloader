package webviewloader

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/juju/mutex"
)

func httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func downloadString(uri string) (string, error) {
	if !strings.HasPrefix(uri, "http") {
		return fileutil.ReadFileToString(uri)
	}

	client := httpClient()
	resp, err := client.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func downloadFile(uri, path string) error {
	if !strings.HasPrefix(uri, "http") {
		fileutil.CreateDir(filepath.Dir(path))
		return fileutil.CopyFile(uri, path)
	}

	client := httpClient()
	resp, err := client.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileutil.CreateDir(filepath.Dir(path))

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func isWindows() bool {
	return strings.EqualFold(runtime.GOOS, "windows")
}

func findName(str string) string {
	str = strings.ReplaceAll(str, "\\", "/")
	arr := strings.Split(str, "/")
	return arr[len(arr)-1]
}

func findBaseName(str string) string {
	name := findName(str)
	dotIndex := strings.Index(name, ".")
	if dotIndex == -1 {
		return name
	}
	return name[:dotIndex]
}

type Clock struct {
}

func (f *Clock) After(t time.Duration) <-chan time.Time {
	return time.After(t)
}
func (f *Clock) Now() time.Time {
	return time.Now()
}
func mutexAcquire(name string, timeout time.Duration) (mutex.Releaser, error) {
	name = replaceMutexName(name)
	spec := mutex.Spec{
		Name:    name,
		Clock:   &Clock{},
		Delay:   time.Millisecond * 300,
		Timeout: timeout,
	}
	return mutex.Acquire(spec)
}

func mutexRelease(releaser mutex.Releaser) {
	defer func() {
		recover()
	}()
	if releaser != nil {
		releaser.Release()
	}
}

func replaceMutexName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	for k := range map[string]struct{}{
		"_": {},
		".": {},
	} {
		name = strings.ReplaceAll(name, k, "-")
	}
	return name
}

type selectURISt struct {
	x64, x64lower, x86, x86lower string
}
