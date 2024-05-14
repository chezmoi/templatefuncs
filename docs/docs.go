package docs

import (
	"embed"
	"log"
	"maps"
	"regexp"
	"strings"
)

type Reference struct {
	Type    string
	Title   string
	Body    string
	Example string
}

//go:embed *.md
var f embed.FS

var (
	References map[string]Reference

	newlineRx   = regexp.MustCompile(`\r?\n`)
	pageTitleRx = regexp.MustCompile(`^#\s+(\S+)`)
	// Template function names must start with a letter or underscore
	// and can subsequently contain letters, underscores and digits.
	funcNameRx = regexp.MustCompile("`" + `([a-zA-Z_]\w*)` + "`")
)

func readFiles() []string {
	fileContents := []string{}

	fileInfos, err := f.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() || !strings.HasSuffix(fileInfo.Name(), ".md") {
			continue
		}
		content, err := f.ReadFile(fileInfo.Name())
		if err != nil {
			log.Fatal(err)
		}
		fileContents = append(fileContents, string(content))
	}

	return fileContents
}

func parseFile(file string) map[string]Reference {
	references := make(map[string]Reference)
	var reference Reference
	var funcName string
	var b strings.Builder
	var e strings.Builder
	inExample := false

	lines := newlineRx.Split(file, -1)
	funcType, lines := pageTitleRx.FindStringSubmatch(lines[0])[1], lines[1:]

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "## "):
			if reference.Title != "" {
				references[funcName] = reference
			}
			funcName = funcNameRx.FindStringSubmatch(line)[1]
			reference = Reference{}
			reference.Type = funcType
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
		references[funcName] = reference
	}

	return references
}

func init() {
	References = make(map[string]Reference)

	files := readFiles()

	for _, file := range files {
		maps.Copy(References, parseFile(file))
	}
}
