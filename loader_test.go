package loadcookie

import (
	"testing"
)

func TestMain(m *testing.M) {
	cfg := Config{
		WebviewName: "webview",
	}
	loader := New(cfg)
	loader.CheckEnv()
}
