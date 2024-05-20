package templatefuncs

import (
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/alecthomas/assert/v2"
)

var (
	strSlice       = []string{"", "a", "b", "c"}
	intSlice       = []int{0, 1, 2, 3}
	mixedSlice     = []any{map[string]any{}, "a", 1, []string{}, 7.7}
	emptySlice     = []any{}
	zeroValueSlice = []int{0, 0, 0, 0}
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
		{
			template: `{{ . | compact }}`,
			data:     strSlice,
			expected: `[a b c]`,
		},
		{
			template: `{{ . | compact }}`,
			data:     intSlice,
			expected: `[1 2 3]`,
		},
		{
			template: `{{ . | compact }}`,
			data:     mixedSlice,
			expected: `[a 1 7.7]`,
		},
		{
			template: `{{ . | compact }}`,
			data:     emptySlice,
			expected: "[]",
		},
		{
			template: `{{ . | compact }}`,
			data:     zeroValueSlice,
			expected: "[]",
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
			template: `{{ list "a" "b" "c" | quote | join "," }}`,
			expected: `"a","b","c"`,
		},
		{
			template: `{{ (lstat "testdata/file").type }}`,
			expected: "file",
		},
		{
			template: `{{ (lstat "testdata/file").type }}`,
			expected: "file",
		},
		{
			template: `{{ prefixLines "# " "a" }}`,
			expected: `# a`,
		},
		{
			template: `{{ prefixLines "# " . }}`,
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
			template: `{{ quote "a" }}`,
			expected: `"a"`,
		},
		{
			template: `{{ "abcba" | replaceAll "b" "d" }}`,
			expected: `adcda`,
		},
		{
			template: `{{ list "abc" "cba" | replaceAll "b" "d" }}`,
			expected: "[adc cda]",
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
			assert.NoError(t, err)
			assert.NotZero(t, tmpl)

			var actual strings.Builder
			assert.NoError(t, tmpl.Execute(&actual, tc.data))
			assert.Equal(t, tc.expected, actual.String())
		})
	}
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}
