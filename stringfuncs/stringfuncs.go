package stringfuncs

import (
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/chezmoi/templatefuncs/internal/transform"
)

func NewFuncMap() template.FuncMap {
	return template.FuncMap{
		"contains":         transform.ReverseArgs2(strings.Contains),
		"hasPrefix":        transform.ReverseArgs2(strings.HasPrefix),
		"hasSuffix":        transform.ReverseArgs2(strings.HasSuffix),
		"quote":            transform.EachString(strconv.Quote),
		"regexpReplaceAll": regexpReplaceAllTemplateFunc,
		"toLower":          transform.EachString(strings.ToLower),
		"toUpper":          transform.EachString(strings.ToUpper),
		"trimSpace":        transform.EachString(strings.TrimSpace),
	}
}

// regexpReplaceAllTemplateFunc is the core implementation of the
// `regexpReplaceAll` template function.
func regexpReplaceAllTemplateFunc(expr, repl, s string) string {
	return regexp.MustCompile(expr).ReplaceAllString(s, repl)
}
