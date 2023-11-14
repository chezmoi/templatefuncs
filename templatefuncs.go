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

var fileModeTypeNames = map[fs.FileMode]string{
	0:                 "file",
	fs.ModeDir:        "dir",
	fs.ModeSymlink:    "symlink",
	fs.ModeNamedPipe:  "named pipe",
	fs.ModeSocket:     "socket",
	fs.ModeDevice:     "device",
	fs.ModeCharDevice: "char device",
}

func NewFuncMap() template.FuncMap {
	return template.FuncMap{
		"comment":          commentTemplateFunc,
		"contains":         reverseArgs2(strings.Contains),
		"eqFold":           eqFoldTemplateFunc,
		"hasPrefix":        reverseArgs2(strings.HasPrefix),
		"hasSuffix":        reverseArgs2(strings.HasSuffix),
		"fromJSON":         eachByteSliceErr(fromJSONTemplateFunc),
		"hexDecode":        eachStringErr(hex.DecodeString),
		"hexEncode":        eachByteSlice(hex.EncodeToString),
		"join":             reverseArgs2(strings.Join),
		"list":             listTemplateFunc,
		"lookPath":         eachStringErr(lookPathTemplateFunc),
		"lstat":            eachString(lstatTemplateFunc),
		"quote":            eachString(strconv.Quote),
		"regexpReplaceAll": regexpReplaceAllTemplateFunc,
		"stat":             eachString(statTemplateFunc),
		"toJSON":           toJSONTemplateFunc,
		"toLower":          eachString(strings.ToLower),
		"toString":         toStringTemplateFunc,
		"toUpper":          eachString(strings.ToUpper),
		"trimSpace":        eachString(strings.TrimSpace),
	}
}

func commentTemplateFunc(prefix, s string) string {
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

func fromJSONTemplateFunc(data []byte) (any, error) {
	var result any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func listTemplateFunc(args ...any) []any {
	return args
}

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

func regexpReplaceAllTemplateFunc(expr, repl, s string) string {
	return regexp.MustCompile(expr).ReplaceAllString(s, repl)
}

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

func toJSONTemplateFunc(arg any) []byte {
	data, err := json.Marshal(arg)
	if err != nil {
		panic(err)
	}
	return data
}

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
		return strconv.FormatFloat(float64(arg), 'f', -1, 64)
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

func reverseArgs2[T1, T2, R any](f func(T1, T2) R) func(T2, T1) R {
	return func(arg1 T2, arg2 T1) R {
		return f(arg2, arg1)
	}
}
