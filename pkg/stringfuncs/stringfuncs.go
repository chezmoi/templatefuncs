package stringfuncs

import (
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/chezmoi/templatefuncs/internal/utils"
)

var FuncMap = template.FuncMap{
	"contains":         utils.ReverseArgs2(strings.Contains),
	"hasPrefix":        utils.ReverseArgs2(strings.HasPrefix),
	"hasSuffix":        utils.ReverseArgs2(strings.HasSuffix),
	"quote":            utils.EachString(strconv.Quote),
	"regexpReplaceAll": regexpReplaceAllTemplateFunc,
	"toLower":          utils.EachString(strings.ToLower),
	"toUpper":          utils.EachString(strings.ToUpper),
	"trimSpace":        utils.EachString(strings.TrimSpace),
}

// regexpReplaceAllTemplateFunc is the core implementation of the
// `regexpReplaceAll` template function.
func regexpReplaceAllTemplateFunc(expr, repl, s string) string {
	return regexp.MustCompile(expr).ReplaceAllString(s, repl)
}
