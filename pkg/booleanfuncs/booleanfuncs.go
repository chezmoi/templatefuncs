package booleanfuncs

import (
	"strings"
	"text/template"
)

var FuncMap = template.FuncMap{
	"eqFold": eqFoldTemplateFunc,
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
