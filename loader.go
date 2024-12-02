package loadcookie

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Config struct {
	WebviewURL             string
	WebviewMd5URL          string
	WebviewWorkDir         string
	DependniesComponentURL string
}

type WebviewOptions struct {
	UA           string
	Title        string
	Width        int
	Height       int
	WaitElements []string
	WaitCookies  []string
}

type WebviewResult struct {
	UA      string            `json:"ua"`
	Cookies map[string]string `json:"cookies"`
}

type Loader struct {
	cfg           Config
	webviewerPath string
}

func New(cfg Config) *Loader {
	return &Loader{cfg: cfg}
}

func (l *Loader) CheckEnv() (err error) {
	err = checkComponent()
	return
}

func (l *Loader) InstallEnv() (err error) {
	err = checkComponent()
	if err != nil {
		err = installComponent(l.cfg.DependniesComponentURL)
	}
	return
}

func (l *Loader) Start(url string, opt WebviewOptions) (result WebviewResult, err error) {
	if l.webviewerPath == "" {
		err = l.CheckEnv()
		if err != nil {
			return
		}
	}

	args := []string{url}
	if opt.Title != "" {
		args = append(args, "--title", opt.Title)
	}
	if opt.Width > 0 {
		args = append(args, "--width", fmt.Sprintf("%d", opt.Width))
	}
	if opt.Height > 0 {
		args = append(args, "--height", fmt.Sprintf("%d", opt.Height))
	}
	if len(opt.WaitElements) > 0 {
		args = append(args, "--elements")
		args = append(args, opt.WaitElements...)
	}
	if len(opt.WaitCookies) > 0 {
		args = append(args, "--cookies")
		args = append(args, opt.WaitCookies...)
	}

	c := exec.Command(l.webviewerPath, args...)
	stdout, err := c.Output()
	if err != nil {
		return
	}
	err = json.Unmarshal(stdout, &result)
	return
}
