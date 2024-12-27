package webviewloader

import (
	"archive/zip"
	"crypto/tls"
	"fmt"
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
	if uri == "" {
		return fmt.Errorf("uri is empty")
	}

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

func findJsonObject(data_ []byte) ([]byte, error) {
	data := string(data_)
	count := 0

	start := 0
	for i, c := range data {
		if c == '{' {
			count++
			if count == 1 {
				start = i
			}
		} else if c == '}' {
			count--
			if count == 0 {
				return []byte(data[start : i+1]), nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

type selectURISt struct {
	x64, x64lower, x86, x86lower string
}

func unZip(zipFile string, destPath string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		path := filepath.Join(destPath, f.Name)

		if !isWindows() {
			isSymlink := f.Mode()&os.ModeSymlink != 0 || (f.ExternalAttrs>>16)&0xA000 == 0xA000

			if isSymlink {
				inFile, err := f.Open()
				if err != nil {
					return fmt.Errorf("failed to open symlink: %v", err)
				}

				linkTarget, err := io.ReadAll(inFile)
				inFile.Close()
				if err != nil {
					return fmt.Errorf("failed to read symlink target: %v", err)
				}

				if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
					return fmt.Errorf("failed to create parent directory for symlink: %v", err)
				}

				if _, err := os.Lstat(path); err == nil {
					os.Remove(path)
				}

				if err := os.Symlink(string(linkTarget), path); err != nil {
					return fmt.Errorf("failed to create symlink %s -> %s: %v", path, string(linkTarget), err)
				}
				continue
			}
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		} else {
			err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				inFile.Close()
				return err
			}

			_, err = io.Copy(outFile, inFile)
			inFile.Close()
			outFile.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
