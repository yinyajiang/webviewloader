package webviewloader

import "fmt"

func checkWebviewComponent() (err error) {
	return nil
}

func installWebviewComponent(...any) (err error) {
	fmt.Print("macos should not install webview")
	return nil
}

func selectURI(st selectURISt) string {
	panic("macos not implement selectURI")
}
