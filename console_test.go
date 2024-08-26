package mintel

import "testing"

func TestConsole(t *testing.T) {
	Register("test", NewConsole(nil, nil, nil))
	client := Open("test", nil)
	client.Logger().Add(KV("LEVEL", "INFO"), KV("Message", "Test")).Flush()
	defer client.Close()
}
