package mintel

import (
	"testing"

	"github.com/vedadiyan/mintel/util/template"
)

func TestConsole(m *testing.T) {
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
	l, _ := template.Parse(logTmlp)
	t, _ := template.Parse(traceTmlp)
	RegisterClient("console", func(m Metadata) Telemetry {
		return NewConsole(l, t, nil)
	})

	tel := Open("console", nil, Trace("Vars", "Ok and Good"))
	defer tel.Close()
}
