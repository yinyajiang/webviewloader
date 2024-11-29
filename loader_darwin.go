package loadcookie

import "fmt"

func checkWebviewEnv(cfg Config) (webviewerPath string, err error) {
	return "", nil
}

func installWebview(cfg Config) (webviewerPath string, err error) {
	fmt.Print("macos should not install webview")
	return "", nil
}
