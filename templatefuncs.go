package templatefuncs

import (
	"errors"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"text/template"

	"github.com/chezmoi/templatefuncs/internal/utils"
	"github.com/chezmoi/templatefuncs/pkg/booleanfuncs"
	"github.com/chezmoi/templatefuncs/pkg/conversionfuncs"
	"github.com/chezmoi/templatefuncs/pkg/listfuncs"
	"github.com/chezmoi/templatefuncs/pkg/stringfuncs"
)

// NewFuncMap returns a new [text/template.FuncMap] containing all template
// functions.
func NewFuncMap() template.FuncMap {
	funcMap := template.FuncMap{}

	maps.Copy(funcMap, template.FuncMap{
		"lookPath": utils.EachStringErr(lookPathTemplateFunc),
		"lstat":    utils.EachString(lstatTemplateFunc),
		"stat":     utils.EachString(statTemplateFunc),
	})

	maps.Copy(funcMap, booleanfuncs.FuncMap)
	maps.Copy(funcMap, conversionfuncs.FuncMap)
	maps.Copy(funcMap, listfuncs.FuncMap)
	maps.Copy(funcMap, stringfuncs.FuncMap)

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
		return utils.FileInfoToMap(fileInfo)
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
		return utils.FileInfoToMap(fileInfo)
	case errors.Is(err, fs.ErrNotExist):
		return nil
	default:
		panic(err)
	}
}
