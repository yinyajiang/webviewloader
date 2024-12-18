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
	if st.x64 != "" {
		return st.x64
	}
	if st.x86 != "" {
		return st.x86
	}
	if st.x64lower != "" {
		return st.x64lower
	}
	if st.x86lower != "" {
		return st.x86lower
	}
	return ""
}
