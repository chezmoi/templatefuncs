package listfuncs

import (
	"strings"
	"text/template"

	"github.com/chezmoi/templatefuncs/internal/utils"
)

var FuncMap = template.FuncMap{
	"join":        utils.ReverseArgs2(strings.Join),
	"list":        listTemplateFunc,
	"prefixLines": prefixLinesTemplateFunc,
}

// listTemplateFunc is the core implementation of the `list` template function.
func listTemplateFunc(args ...any) []any {
	return args
}

// prefixLinesTemplateFunc is the core implementation of the `prefixLines`
// template function.
func prefixLinesTemplateFunc(prefix, s string) string {
	type stateType int
	const (
		startOfLine stateType = iota
		inLine
	)

	state := startOfLine
	var builder strings.Builder
	builder.Grow(2 * len(s))
	for _, r := range s {
		switch state {
		case startOfLine:
			if _, err := builder.WriteString(prefix); err != nil {
				panic(err)
			}
			if _, err := builder.WriteRune(r); err != nil {
				panic(err)
			}
			if r != '\n' {
				state = inLine
			}
		case inLine:
			if _, err := builder.WriteRune(r); err != nil {
				panic(err)
			}
			if r == '\n' {
				state = startOfLine
			}
		}
	}
	return builder.String()
}
