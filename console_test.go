package mintel

import (
	"net/http"
	"testing"
)

func TestConsole(t *testing.T) {
	Register("test", NewConsole(nil, nil, nil))

	fn := func() {
		testTace, _ := http.NewRequest(http.MethodGet, "https://www.google.com", nil)

		client := Open("test", nil, TraceRef("req", testTace.URL), Trace("method", testTace.Method))
		client.Logger().Add(KV("LEVEL", "INFO"), KV("Message", "Test")).Flush()
		defer client.Close()
	}

	fn()
	fn()
}
