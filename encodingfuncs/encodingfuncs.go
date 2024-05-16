package encodingfuncs

import (
	"encoding/hex"
	"encoding/json"
	"text/template"

	"github.com/chezmoi/templatefuncs/internal/transform"
)

func NewFuncMap() template.FuncMap {
	return template.FuncMap{
		"fromJSON":  transform.EachByteSliceErr(fromJSONTemplateFunc),
		"hexDecode": transform.EachStringErr(hex.DecodeString),
		"hexEncode": transform.EachByteSlice(hex.EncodeToString),
		"toJSON":    toJSONTemplateFunc,
	}
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
