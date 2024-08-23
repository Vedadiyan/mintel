package template

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"text/template"
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
)

var (
	mainPattern  = regexp.MustCompile(`\$\([^\(^\)]*\)`)
	indexPattern = regexp.MustCompile(`\[[^\[^\]]*\]`)
	digit        = regexp.MustCompile(`[0-9]+$`)
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
	tw.segments = segments
	tw.length = len(segments)
	ended := false
	offset := 0
	if tw.fromMap {
		offset = 1
		ended = tw.write(fmt.Sprintf("[%s]", tw.segments[0]))
	}
	for _, seg := range tw.segments[offset:] {
		if tw.write(seg) {
			ended = true
			break
		}
	}
	if !ended {
		tw.left.WriteString("{{ . }}")
	}
	b.Write(tw.left.Bytes())
	b.Write(tw.right.Bytes())
}

func (tw *TemplateWriter) wildCard(seg string) bool {
	isLast := tw.index == tw.length-1
	if isLast {
		tw.left.WriteString("{ ")
	}
	tw.left.WriteString("{{- range $i, $v := . }}")
	tw.left.WriteString("{{- with $v }}")
	tw.right.WriteString("{{- end }}")
	if isLast {
		tw.left.WriteString("{{- if gt $i 0 }}, {{- end}}\"{{ $i }}\": \"{{ . }}\"")
		tw.left.WriteString("{{- end }}")
		tw.left.WriteString(" }")
		return true
	}
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

func Generate(str string) (map[string]*template.Template, error) {
	matches := mainPattern.FindAllString(str, -1)
	out := make(map[string]*template.Template)
	tw := New(TreatTopAsMap())
	for _, match := range matches {
		var buffer bytes.Buffer
		tw.Write(match, &buffer)
		value := buffer.String()
		fmt.Println(value)
		template, err := template.New(value).Parse(value)
		if err != nil {
			return nil, err
		}
		out[match] = template
	}
	return out, nil
}
