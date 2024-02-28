package docs

import (
	_ "embed"
	"regexp"
	"strings"
)

//go:embed templatefuncs.md
var templateFuncsStr string

type Reference struct {
	Title   string
	Body    string
	Example string
}

var References map[string]Reference

func init() {
	newlineRx := regexp.MustCompile(`\r?\n`)

	// Template function names must start with a letter or underscore
	// and can subsequently contain letters, underscores and digits.
	funcNameRx := regexp.MustCompile("`" + `([a-zA-Z_]\w*)` + "`")

	References = make(map[string]Reference)
	var reference Reference
	var funcName string
	var b strings.Builder
	var e strings.Builder
	inExample := false

	for _, line := range newlineRx.Split(templateFuncsStr, -1) {
		switch {
		case strings.HasPrefix(line, "## "):
			if reference.Title != "" {
				References[funcName] = reference
			}
			funcName = funcNameRx.FindStringSubmatch(line)[1]
			reference = Reference{}
			reference.Title = strings.TrimPrefix(line, "## ")
		case strings.HasPrefix(line, "```"):
			if !inExample {
				reference.Body = strings.TrimSpace(b.String())
				b.Reset()
			}
			e.WriteString(line + "\n")
			if inExample {
				reference.Example = strings.TrimSpace(e.String())
				e.Reset()
			}
			inExample = !inExample
		case inExample:
			if _, err := e.WriteString(line + "\n"); err != nil {
				panic(err)
			}
		case reference.Title != "":
			if _, err := b.WriteString(line + "\n"); err != nil {
				panic(err)
			}
		}
	}

	if reference.Title != "" {
		References[funcName] = reference
	}
}
