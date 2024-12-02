package loadcookie

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/duke-git/lancet/v2/fileutil"
)

func httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func downloadFile(url, path string) error {
	client := httpClient()
	resp, err := client.Get(url)
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
