package conversionfuncs

import (
	"fmt"
	"strconv"
	"text/template"
)

func NewFuncMap() template.FuncMap {
	return template.FuncMap{
		"toString": toStringTemplateFunc,
	}
}

// toStringTemplateFunc is the core implementation of the `toString` template
// function.
func toStringTemplateFunc(arg any) string {
	// FIXME add more types
	switch arg := arg.(type) {
	case string:
		return arg
	case []byte:
		return string(arg)
	case bool:
		return strconv.FormatBool(arg)
	case float32:
		return strconv.FormatFloat(float64(arg), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(arg, 'f', -1, 64)
	case int:
		return strconv.Itoa(arg)
	case int32:
		return strconv.FormatInt(int64(arg), 10)
	case int64:
		return strconv.FormatInt(arg, 10)
	default:
		panic(fmt.Sprintf("%T: unsupported type", arg))
	}
}
