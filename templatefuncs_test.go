package templatefuncs

import (
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEachString(t *testing.T) {
	for i, tc := range []struct {
		f        func(string) string
		arg      any
		expected any
	}{
		{
			f:        strings.ToUpper,
			arg:      "a",
			expected: "A",
		},
		{
			f:        strings.ToUpper,
			arg:      []string{"a", "b", "c"},
			expected: []string{"A", "B", "C"},
		},
		{
			f:        strings.ToUpper,
			arg:      []byte("a"),
			expected: "A",
		},
		{
			f:        strings.ToUpper,
			arg:      [][]byte{[]byte("a"), []byte("b"), []byte("c")},
			expected: []string{"A", "B", "C"},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			f := eachString(tc.f)
			assert.Equal(t, tc.expected, f(tc.arg))
		})
	}
}

func TestFuncMap(t *testing.T) {
	funcMap := NewFuncMap()
	for i, tc := range []struct {
		template string
		data     any
		expected string
	}{
		{},
		{
			template: `{{ comment "# " "a" }}`,
			expected: `# a`,
		},
		{
			template: `{{ comment "# " . }}`,
			data: joinLines(
				"a",
				"b",
			),
			expected: joinLines(
				"# a",
				"# b",
			),
		},
		{
			template: `{{ "abc" | contains "bc" }}`,
			expected: "true",
		},
		{
			template: `{{ "abc" | contains "cd" }}`,
			expected: "false",
		},
		{
			template: `{{ eqFold "A" "a" }}`,
			expected: "true",
		},
		{
			template: `{{ eqFold "B" "a" "b" }}`,
			expected: "true",
		},
		{
			template: `{{ eqFold "C" "a" "b" }}`,
			expected: "false",
		},
		{
			template: `{{ fromJSON "0" }}`,
			expected: "0",
		},
		{
			template: `{{ "ab" | hasPrefix "a" }}`,
			expected: "true",
		},
		{
			template: `{{ "ab" | hasPrefix "b" }}`,
			expected: "false",
		},
		{
			template: `{{ "ab" | hasSuffix "a" }}`,
			expected: "false",
		},
		{
			template: `{{ "ab" | hasSuffix "b" }}`,
			expected: "true",
		},
		{
			template: `{{ (lstat "testdata/file").type }}`,
			expected: "file",
		},
		{
			template: `{{ quote "a" }}`,
			expected: `"a"`,
		},
		{
			template: `{{ list "a" "b" "c" | quote | join "," }}`,
			expected: `"a","b","c"`,
		},
		{
			template: `{{ (lstat "testdata/file").type }}`,
			expected: "file",
		},
		{
			template: `{{ (stat "testdata/file").type }}`,
			expected: "file",
		},
		{
			template: `{{ toJSON 0 | toString }}`,
			expected: "0",
		},
		{
			template: `{{ trimSpace " a " }}`,
			expected: "a",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tmpl, err := template.New("").Funcs(funcMap).Parse(tc.template)
			require.NoError(t, err)
			require.NotNil(t, tmpl)

			var actual strings.Builder
			require.NoError(t, tmpl.Execute(&actual, tc.data))
			assert.Equal(t, tc.expected, actual.String())
		})
	}
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}
