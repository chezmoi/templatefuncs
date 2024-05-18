package templatefuncs

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// fileModeTypeNames maps file mode types to human-readable strings.
var fileModeTypeNames = map[fs.FileMode]string{
	0:                 "file",
	fs.ModeDir:        "dir",
	fs.ModeSymlink:    "symlink",
	fs.ModeNamedPipe:  "named pipe",
	fs.ModeSocket:     "socket",
	fs.ModeDevice:     "device",
	fs.ModeCharDevice: "char device",
}

// NewFuncMap returns a new [text/template.FuncMap] containing all template
// functions.
func NewFuncMap() template.FuncMap {
	return template.FuncMap{
		"contains":         reverseArgs2(strings.Contains),
		"eqFold":           eqFoldTemplateFunc,
		"fromJSON":         eachByteSliceErr(fromJSONTemplateFunc),
		"hasPrefix":        reverseArgs2(strings.HasPrefix),
		"hasSuffix":        reverseArgs2(strings.HasSuffix),
		"hexDecode":        eachStringErr(hex.DecodeString),
		"hexEncode":        eachByteSlice(hex.EncodeToString),
		"join":             reverseArgs2(strings.Join),
		"list":             listTemplateFunc,
		"lookPath":         eachStringErr(lookPathTemplateFunc),
		"lstat":            eachString(lstatTemplateFunc),
		"prefixLines":      prefixLinesTemplateFunc,
		"quote":            eachString(strconv.Quote),
		"regexpReplaceAll": regexpReplaceAllTemplateFunc,
		"replaceAll":       replaceAllTemplateFunc,
		"stat":             eachString(statTemplateFunc),
		"toJSON":           toJSONTemplateFunc,
		"toLower":          eachString(strings.ToLower),
		"toString":         toStringTemplateFunc,
		"toUpper":          eachString(strings.ToUpper),
		"trimSpace":        eachString(strings.TrimSpace),
	}
}

// eqFoldTemplateFunc is the core implementation of the `eqFold` template
// function.
func eqFoldTemplateFunc(first, second string, more ...string) bool {
	if strings.EqualFold(first, second) {
		return true
	}
	for _, s := range more {
		if strings.EqualFold(first, s) {
			return true
		}
	}
	return false
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

// listTemplateFunc is the core implementation of the `list` template function.
func listTemplateFunc(args ...any) []any {
	return args
}

// lookPathTemplateFunc is the core implementation of the `lookPath` template
// function.
func lookPathTemplateFunc(file string) (string, error) {
	switch path, err := exec.LookPath(file); {
	case err == nil:
		return path, nil
	case errors.Is(err, exec.ErrNotFound):
		return "", nil
	case errors.Is(err, fs.ErrNotExist):
		return "", nil
	default:
		return "", err
	}
}

// lstatTemplateFunc is the core implementation of the `lstat` template
// function.
func lstatTemplateFunc(name string) any {
	switch fileInfo, err := os.Lstat(name); {
	case err == nil:
		return fileInfoToMap(fileInfo)
	case errors.Is(err, fs.ErrNotExist):
		return nil
	default:
		panic(err)
	}
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

// replaceAllTemplateFunc is the `replaceAll` template function.
func replaceAllTemplateFunc(old, new string, v any) any { //nolint:predeclared
	return eachString(func(s string) any {
		return strings.ReplaceAll(s, old, new)
	})(v)
}

// regexpReplaceAllTemplateFunc is the core implementation of the
// `regexpReplaceAll` template function.
func regexpReplaceAllTemplateFunc(expr, repl, s string) string {
	return regexp.MustCompile(expr).ReplaceAllString(s, repl)
}

// statTemplateFunc is the core implementation of the `stat` template function.
func statTemplateFunc(name string) any {
	switch fileInfo, err := os.Stat(name); {
	case err == nil:
		return fileInfoToMap(fileInfo)
	case errors.Is(err, fs.ErrNotExist):
		return nil
	default:
		panic(err)
	}
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

// eachByteSlice transforms a function that takes a single `[]byte` and returns
// a `T` to a function that takes zero or more `[]byte`-like arguments and
// returns zero or more `T`s.
func eachByteSlice[T any](f func([]byte) T) func(any) any {
	return func(arg any) any {
		switch arg := arg.(type) {
		case []byte:
			return f(arg)
		case [][]byte:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				result = append(result, f(a))
			}
			return result
		case string:
			return f([]byte(arg))
		case []string:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				result = append(result, f([]byte(a)))
			}
			return result
		default:
			panic(fmt.Sprintf("%T: unsupported argument type", arg))
		}
	}
}

// eachByteSliceErr transforms a function that takes a single `[]byte` and
// returns a `T` and an `error` into a function that takes zero or more
// `[]byte`-like arguments and returns zero or more `Ts` and an error.
func eachByteSliceErr[T any](f func([]byte) (T, error)) func(any) any {
	return func(arg any) any {
		switch arg := arg.(type) {
		case []byte:
			result, err := f(arg)
			if err != nil {
				panic(err)
			}
			return result
		case [][]byte:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				r, err := f(a)
				if err != nil {
					panic(err)
				}
				result = append(result, r)
			}
			return result
		case string:
			result, err := f([]byte(arg))
			if err != nil {
				panic(err)
			}
			return result
		case []string:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				r, err := f([]byte(a))
				if err != nil {
					panic(err)
				}
				result = append(result, r)
			}
			return result
		default:
			panic(fmt.Sprintf("%T: unsupported argument type", arg))
		}
	}
}

// eachString transforms a function that takes a single `string`-like argument
// and returns a `T` into a function that takes zero or more `string`-like
// arguments and returns zero or more `T`s.
func eachString[T any](f func(string) T) func(any) any {
	return func(arg any) any {
		switch arg := arg.(type) {
		case string:
			return f(arg)
		case []string:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				result = append(result, f(a))
			}
			return result
		case []byte:
			return f(string(arg))
		case [][]byte:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				result = append(result, f(string(a)))
			}
			return result
		case []any:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				switch a := a.(type) {
				case string:
					result = append(result, f(a))
				case []byte:
					result = append(result, f(string(a)))
				default:
					panic(fmt.Sprintf("%T: unsupported argument type", a))
				}
			}
			return result
		default:
			panic(fmt.Sprintf("%T: unsupported argument type", arg))
		}
	}
}

// eachStringErr transforms a function that takes a single `string`-like argument
// and returns a `T` and an `error` into a function that takes zero or more
// `string`-like arguments and returns zero or more `T`s and an `error`.
func eachStringErr[T any](f func(string) (T, error)) func(any) any {
	return func(arg any) any {
		switch arg := arg.(type) {
		case string:
			result, err := f(arg)
			if err != nil {
				panic(err)
			}
			return result
		case []string:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				r, err := f(a)
				if err != nil {
					panic(err)
				}
				result = append(result, r)
			}
			return result
		case []byte:
			result, err := f(string(arg))
			if err != nil {
				panic(err)
			}
			return result
		case [][]byte:
			result := make([]T, 0, len(arg))
			for _, a := range arg {
				r, err := f(string(a))
				if err != nil {
					panic(err)
				}
				result = append(result, r)
			}
			return result
		default:
			panic(fmt.Sprintf("%T: unsupported argument type", arg))
		}
	}
}

// fileInfoToMap returns a `map[string]any` of `fileInfo`'s fields.
func fileInfoToMap(fileInfo fs.FileInfo) map[string]any {
	return map[string]any{
		"name":    fileInfo.Name(),
		"size":    fileInfo.Size(),
		"mode":    int(fileInfo.Mode()),
		"perm":    int(fileInfo.Mode().Perm()),
		"modTime": fileInfo.ModTime().Unix(),
		"isDir":   fileInfo.IsDir(),
		"type":    fileModeTypeNames[fileInfo.Mode()&fs.ModeType],
	}
}

// reverseArgs2 transforms a function that takes two arguments and returns an
// `R` into a function that takes the arguments in reverse order and returns an
// `R`.
func reverseArgs2[T1, T2, R any](f func(T1, T2) R) func(T2, T1) R {
	return func(arg1 T2, arg2 T1) R {
		return f(arg2, arg1)
	}
}
