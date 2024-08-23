package template

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/vedadiyan/mintel/util/json"
)

type (
	TemplateWriter struct {
		length   int
		index    int
		segments []string
		fromMap  bool

		left  bytes.Buffer
		right bytes.Buffer

		mut sync.Mutex
	}
	TemplateWriterOpts func(*TemplateWriter)
	BinderFunc         func(v any) string
	Binder             map[string]BinderFunc
)

var (
	mainPattern  = regexp.MustCompile(`\$\([^\(^\)]*\)`)
	indexPattern = regexp.MustCompile(`\[[^\[^\]]*\]`)
	digit        = regexp.MustCompile(`[0-9]+$`)
	space        = regexp.MustCompile(`(("[^"]*"|[^"\s]*)*)\s*`)
)

func TreatTopAsMap() TemplateWriterOpts {
	return func(tw *TemplateWriter) {
		tw.fromMap = true
	}
}

func New(opts ...TemplateWriterOpts) *TemplateWriter {
	tw := new(TemplateWriter)
	for _, opt := range opts {
		opt(tw)
	}
	return tw
}

func (tw *TemplateWriter) Write(str string, b io.Writer) {
	tw.mut.Lock()
	defer tw.mut.Unlock()
	defer tw.left.Reset()
	defer tw.right.Reset()
	str = strings.TrimLeftFunc(str, func(r rune) bool { return r == '$' || r == '(' || r == ' ' })
	str = strings.TrimRightFunc(str, func(r rune) bool { return r == ')' })
	segments := strings.Split(str, ".")
	tw.index = 0
	tw.segments = segments
	tw.length = len(segments)
	ended := false
	offset := 0
	if tw.fromMap {
		offset = 1
		ended = tw.write(fmt.Sprintf("[%s]", tw.segments[0]))
	}
	for index, seg := range tw.segments[offset:] {
		tw.index = index + offset
		if tw.write(seg) {
			ended = true
			break
		}
	}
	if !ended {
		tw.left.WriteString("{{ Serialize . }}")
	}
	b.Write(tw.left.Bytes())
	b.Write(tw.right.Bytes())
}

func (tw *TemplateWriter) wildCard(seg string) bool {
	isLast := tw.index == tw.length-1
	if isLast {
		tw.left.WriteString("{{- Serialize . }}")
		return true
	}
	tw.left.WriteString("{{- range $i, $v := . }}")
	tw.left.WriteString("{{- with $v }}")
	tw.right.WriteString("{{- end }}")
	tw.right.WriteString("{{- end }}")
	return false
}

func (tw *TemplateWriter) field(seg string) bool {
	tw.right.WriteString("{{- end }}")
	tw.left.WriteString("{{- with .")
	tw.left.WriteString(seg)
	tw.left.WriteString(" }}")
	return false
}

func (tw *TemplateWriter) key(seg string) bool {
	segments := strings.Split(seg, "[")
	if len(segments) == 1 {
		tw.right.WriteString("{{- end }}")
		tw.left.WriteString("{{- with .")
		tw.left.WriteString(segments[0])
		tw.left.WriteString(" }}")
		return false
	}
	for _, seg := range segments {
		if len(seg) == 0 {
			continue
		}
		tw.right.WriteString("{{- end }}")
		seg = strings.TrimLeftFunc(seg, func(r rune) bool { return r == '[' || r == '"' })
		seg = strings.TrimRightFunc(seg, func(r rune) bool { return r == ']' || r == '"' })
		isArrayIndex := digit.MatchString(seg)
		tw.left.WriteString("{{- with index . ")
		if !isArrayIndex {
			tw.left.WriteString("\"")
		}
		tw.left.WriteString(seg)
		if !isArrayIndex {
			tw.left.WriteString("\"")
		}
		tw.left.WriteString("}}")
	}
	return false
}

func (tw *TemplateWriter) write(seg string) bool {
	if seg == "*" {
		return tw.wildCard(seg)
	}
	if !indexPattern.MatchString(seg) {
		return tw.field(seg)
	}
	return tw.key(seg)
}

func Parse(templateStr string) (Binder, error) {
	templateStr = strings.ReplaceAll(templateStr, "\r", "")
	templateStr = strings.ReplaceAll(templateStr, "\n", "")
	templateStr = strings.TrimLeftFunc(templateStr, func(r rune) bool { return r == ' ' || r == '\t' })
	templateStr = strings.TrimRightFunc(templateStr, func(r rune) bool { return r == ' ' || r == '\t' })
	templateStr = RemoveSpace(templateStr)
	matches := mainPattern.FindAllString(templateStr, -1)
	out := make(Binder)
	out["_"] = func(v any) string {
		return templateStr
	}
	tw := New(TreatTopAsMap())
	serialize := func(v any) string {
		return string(json.Marshal(v))
	}
	for _, match := range matches {
		var buffer bytes.Buffer
		tw.Write(match, &buffer)
		value := buffer.String()
		//fmt.Println(value)
		template, err := template.New(value).Funcs(template.FuncMap{
			"Serialize": serialize,
		}).Parse(value)
		if err != nil {
			return nil, err
		}
		out[match] = func(v any) string {
			var buffer bytes.Buffer
			_ = template.Execute(&buffer, v)
			return buffer.String()
		}
	}
	return out, nil
}

func Bind(binder Binder, v any) string {
	templateStr := binder["_"](nil)
	for key, value := range binder {
		value := value(v)
		if len(value) == 0 {
			value = "null"
		}
		templateStr = strings.ReplaceAll(templateStr, key, value)
	}
	return templateStr
}

func RemoveSpace(str string) string {
	var buffer bytes.Buffer
	hold := false
	for i := 0; i < len(str); i++ {
		r := str[i]
		switch r {
		case '\\':
			{
				i++
				continue
			}
		case '"':
			{
				hold = !hold
			}
		}
		if hold {
			buffer.WriteByte(r)
			continue
		}
		if r != ' ' && r != '\t' {
			buffer.WriteByte(r)
		}
	}
	return buffer.String()
}
