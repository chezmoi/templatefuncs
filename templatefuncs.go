package templatefuncs

import (
	"errors"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"text/template"

	"github.com/chezmoi/templatefuncs/booleanfuncs"
	"github.com/chezmoi/templatefuncs/conversionfuncs"
	"github.com/chezmoi/templatefuncs/encodingfuncs"
	"github.com/chezmoi/templatefuncs/internal/transform"
	"github.com/chezmoi/templatefuncs/listfuncs"
	"github.com/chezmoi/templatefuncs/stringfuncs"
)

// NewFuncMap returns a new [text/template.FuncMap] containing all template
// functions.
func NewFuncMap() template.FuncMap {
	funcMap := make(template.FuncMap)

	maps.Copy(funcMap, template.FuncMap{
		"lookPath": transform.EachStringErr(lookPathTemplateFunc),
		"lstat":    transform.EachString(lstatTemplateFunc),
		"stat":     transform.EachString(statTemplateFunc),
	})

	maps.Copy(funcMap, booleanfuncs.NewFuncMap())
	maps.Copy(funcMap, conversionfuncs.NewFuncMap())
	maps.Copy(funcMap, encodingfuncs.NewFuncMap())
	maps.Copy(funcMap, listfuncs.NewFuncMap())
	maps.Copy(funcMap, stringfuncs.NewFuncMap())

	return funcMap
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
		return transform.FileInfoToMap(fileInfo)
	case errors.Is(err, fs.ErrNotExist):
		return nil
	default:
		panic(err)
	}
}

// statTemplateFunc is the core implementation of the `stat` template function.
func statTemplateFunc(name string) any {
	switch fileInfo, err := os.Stat(name); {
	case err == nil:
		return transform.FileInfoToMap(fileInfo)
	case errors.Is(err, fs.ErrNotExist):
		return nil
	default:
		panic(err)
	}
}
