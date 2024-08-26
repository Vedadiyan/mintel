package mintel

import (
	"fmt"
	"net/url"
	"testing"
)

func TestConsole(t *testing.T) {
	Register("test", NewConsole(nil, nil, nil))

	fn := func() {
		testTace, _ := url.Parse("https://www.google.com")
		var err error
		client := Open("test", nil, TraceRef("req", &testTace), Trace("method", "GET"))
		defer client.Close()
		client.Logger().Add(Begin(), Timestamp()).Flush()
		defer client.Logger().Add(End(), Timestamp())
		defer func() {
			if err == nil {
				client.Meter().Add(KV("Call", 1))
				return
			}
			client.Logger().Add(Error(), KV("message", err.Error()), Timestamp())
		}()
		err = fmt.Errorf("test error")

		// testTace = &url.URL{}
		client.Logger().Add(Info(), KV("message", "Test")).Flush()
	}

	fn()
	fn()
}
