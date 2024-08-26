package mintel

import (
	"testing"
)

func TestConsole(t *testing.T) {
	Register("test", NewConsole(nil, nil, nil))

	fn := func() {
		// testTace, _ := url.Parse("https://www.google.com")
		// var err error
		client := Open("test", nil)
		// defer client.Close()
		// defer func() {
		// 	if err == nil {
		// 		client.Meter().Add(KV("Call", 1))
		// 		return
		// 	}
		// 	client.Logger().Add(Error(), KV("Message", err.Error()), Timestamp())
		// }()
		// err = fmt.Errorf("test error")

		// testTace = &url.URL{}
		client.Logger().Add(Info(), KV("message", "Test"), Timestamp()).Flush()
	}

	fn()
	fn()
}
