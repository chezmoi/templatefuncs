package conversionfuncs

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"text/template"

	"github.com/chezmoi/templatefuncs/internal/utils"
)

var FuncMap = template.FuncMap{
	"fromJSON":  utils.EachByteSliceErr(fromJSONTemplateFunc),
	"hexDecode": utils.EachStringErr(hex.DecodeString),
	"hexEncode": utils.EachByteSlice(hex.EncodeToString),
	"toJSON":    toJSONTemplateFunc,
	"toString":  toStringTemplateFunc,
}

// fromJSONTemplateFunc is the core implementation of the `fromJSON` template
// function.
func fromJSONTemplateFunc(data []byte) (any, error) {
	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// toJSONTemplateFunc is the core implementation of the `toJSON` template
// function.
func toJSONTemplateFunc(arg any) []byte {
	data, err := json.Marshal(arg)
	if err != nil {
		panic(err)
	}
	return data
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
