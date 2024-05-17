package transform

import (
	"fmt"
	"io/fs"
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

// EachByteSlice transforms a function that takes a single `[]byte` and returns
// a `T` to a function that takes zero or more `[]byte`-like arguments and
// returns zero or more `T`s.
func EachByteSlice[T any](f func([]byte) T) func(any) any {
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

// EachByteSliceErr transforms a function that takes a single `[]byte` and
// returns a `T` and an `error` into a function that takes zero or more
// `[]byte`-like arguments and returns zero or more `Ts` and an error.
func EachByteSliceErr[T any](f func([]byte) (T, error)) func(any) any {
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

// EachString transforms a function that takes a single `string`-like argument
// and returns a `T` into a function that takes zero or more `string`-like
// arguments and returns zero or more `T`s.
func EachString[T any](f func(string) T) func(any) any {
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

// EachStringErr transforms a function that takes a single `string`-like argument
// and returns a `T` and an `error` into a function that takes zero or more
// `string`-like arguments and returns zero or more `T`s and an `error`.
func EachStringErr[T any](f func(string) (T, error)) func(any) any {
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

// FileInfoToMap returns a `map[string]any` of `fileInfo`'s fields.
func FileInfoToMap(fileInfo fs.FileInfo) map[string]any {
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

// ReverseArgs2 transforms a function that takes two arguments and returns an
// `R` into a function that takes the arguments in reverse order and returns an
// `R`.
func ReverseArgs2[T1, T2, R any](f func(T1, T2) R) func(T2, T1) R {
	return func(arg1 T2, arg2 T1) R {
		return f(arg2, arg1)
	}
}
