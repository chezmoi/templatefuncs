package docs_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/chezmoi/templatefuncs/docs"
)

func TestReferences(t *testing.T) {
	assert.Equal(t, docs.Reference{
		Type:  "String",
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
		Type:  "Boolean",
		Title: "`eqFold` *string1* *string2* [*extraStrings*...]",
		Body: "" +
			"`eqFold` returns the boolean truth of comparing *string1* with *string2*\n" +
			"and any number of *extraStrings* under Unicode case-folding.",
		Example: "" +
			"```text\n" +
			"{{ eqFold \"föö\" \"FOO\" }}\n" +
			"\n" +
			"true\n" +
			"```",
	}, docs.References["eqFold"])
}
