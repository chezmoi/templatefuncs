package docs_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/chezmoi/templatefuncs/docs"
)

func TestReferences(t *testing.T) {
	assert.Equal(t, docs.Reference{
		Title: "`contains` *substring* *string*",
		Body:  "`contains` returns whether *substring* is in *string*.",
		Example: "" +
			"```text\n" +
			"{{ \"abc\" | contains \"ab\" }}\n" +
			"\n" +
			"true\n" +
			"```",
	}, docs.References["contains"])
	assert.Equal(t, docs.Reference{
		Title: "`trimSpace` *string*",
		Body:  "`trimSpace` returns *string* with all spaces removed.",
		Example: "```text\n" +
			"{{ \"    foobar    \" | trimSpace }}\n" +
			"\n" +
			"foobar\n" +
			"```",
	}, docs.References["trimSpace"])
}
