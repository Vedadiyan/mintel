package mintel

import "testing"

func TestConsole(m *testing.T) {
	RegisterClient("console", func(m Metadata) Telemetry {
		logTmlp := `
{
		"app":  $(App),
		"func": $(Func) 
}
		`

		traceTmlp := `
		{
				"vars":  $(Vars),
				"args":  $(Args) 
		}
				`

		v, _ := NewConsole(logTmlp, traceTmlp, "")
		return v
	})

	t := Open("console", nil, Trace("Vars", "Ok and Good"))
	defer t.Close()
}
