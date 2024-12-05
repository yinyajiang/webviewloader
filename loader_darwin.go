package webviewloader

import "fmt"

func checkComponent() (err error) {
	return nil
}

func installComponent(...any) (err error) {
	fmt.Print("macos should not install webview")
	return nil
}

func selectURI(st selectURISt) string {
	panic("macos not implement selectURI")
}
