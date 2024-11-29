package loadcookie

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Config struct {
	WebviewerURI     string
	WebviewerMd5URI  string
	WebviewerWorkDir string
	ComposURI        string

	Title        string
	Width        int
	Height       int
	WaitElements []string
	WaitCookies  []string
}

type ResultInfo struct {
	UA      string            `json:"ua"`
	Cookies map[string]string `json:"cookies"`
}

func Start(url string, cfg Config) (result ResultInfo, err error) {
	webviewerPath, err := checkWebviewEnv(cfg)
	if err != nil {
		webviewerPath, err = installWebview(cfg)
		if err != nil {
			return
		}
	}

	args := []string{
		"--url", url,
	}
	if cfg.Title != "" {
		args = append(args, "--title", cfg.Title)
	}
	if cfg.Width > 0 {
		args = append(args, "--width", fmt.Sprintf("%d", cfg.Width))
	}
	if cfg.Height > 0 {
		args = append(args, "--height", fmt.Sprintf("%d", cfg.Height))
	}
	if len(cfg.WaitElements) > 0 {
		args = append(args, "--elements")
		args = append(args, cfg.WaitElements...)
	}
	if len(cfg.WaitCookies) > 0 {
		args = append(args, "--cookies")
		args = append(args, cfg.WaitCookies...)
	}

	c := exec.Command(webviewerPath, args...)
	stdout, err := c.Output()
	if err != nil {
		return
	}
	err = json.Unmarshal(stdout, &result)
	return
}
