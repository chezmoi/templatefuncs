# Template Functions

## `compact` *list*

`compact` removes all zero value items from *list*.

```text
{{ list "one" "" list "three" | compact }}

[one three]
```

## `contains` *substring* *string*

`contains` returns whether *substring* is in *string*.

```text
{{ "abc" | contains "ab" }}

true
```

## `eqFold` *string1* *string2* [*extraStrings*...]

`eqFold` returns the boolean truth of comparing *string1* with *string2*
and any number of *extraStrings* under Unicode case-folding.

```text
{{ eqFold "föö" "FOO" }}

true
```

## `fromJSON` *jsontext*

`fromJSON` parses *jsontext* as JSON and returns the parsed value.

```text
{{ `{ "foo": "bar" }` | fromJSON }}
```

## `has` *item* *list*

`has` returns whether *item* is in *list*.

```text
{{ list 1 2 3 | has 3 }}

true
```

## `hasPrefix` *prefix* *string*

`hasPrefix` returns whether *string* begins with *prefix*.

```text
{{ "foobar" | hasPrefix "foo" }}

true
```

## `hasSuffix` *suffix* *string*

`hasSuffix` returns whether *string* ends with *suffix*.

```text
{{ "foobar" | hasSuffix "bar" }}

true
```

## `hexDecode` *hextext*

`hexDecode` returns the bytes represented by *hextext*.

```text
{{ hexDecode "666f6f626172" }}

foobar
```

## `hexEncode` *string*

`hexEncode` returns the hexadecimal encoding of *string*.

```text
{{ hexEncode "foobar" }}

666f6f626172
```

## `indexOf` *item* *list*

`indexOf` returns the index of *item* in *list*, or -1 if *item* is not
in *list*.

```text
{{ list "a" "b" "c" | indexOf "b" }}

1
```

## `join` *delimiter* *list*

`join` returns a string containing each item in *list* joined with *delimiter*.

```text
{{ list "foo" "bar" "baz" | join "," }}

foo,bar,baz
```

## `list` *items*...

`list` creates a new list containing *items*.

```text
{{ list "foo" "bar" "baz" }}
```

## `lookPath` *file*

`lookPath` searches for the executable *file* in the users `PATH`
environment variable and returns its path.

```text
{{ lookPath "git" }}
```

## `lstat` *path*

`lstat` returns a map representation of executing
[`os.Lstat`](https://pkg.go.dev/os#Lstat) on *path*.

```text
{{ (lstat "some/file").type }}

file
```

## `prefixLines` *prefix* *list*

`prefixLines` returns a string consisting of each item in *list*
with the prefix *prefix*.

```text
{{ list "this is" "a multi-line" "comment" | prefixLines "# " }}

# this is
# a multi-line
# comment
```

## `quote` *input*

`quote` returns a double-quoted string literal containing *input*.
*input* can be a string or list of strings.

```text
{{ "foobar" | quote }}

"foobar"
```

## `regexpReplaceAll` *pattern* *replacement* *string*

`regexpReplaceAll` replaces all instances of *pattern*
with *replacement* in *string*.

```text
{{ "foobar" | regexpReplaceAll "o*b" "" }}

far
```

## `replaceAll` *old* *new* *string*

`replaceAll` replaces all instances of *old* with *new* in *string*.

```text
{{ "abcba" | replaceAll "b" "d" }}

adcda
```

## `stat` *path*

`stat` returns a map representation of executing
[`os.Stat`](https://pkg.go.dev/os#Stat) on *path*.

```text
{{ (stat "some/file").type }}

file
```

## `toJSON` *input*

`toJSON` returns a JSON string representation of *input*.

```text
{{ list "foo" "bar" "baz" }}

["foo","bar","baz"]
```

## `toLower` *string*

`toLower` returns *string* with all letters converted to lower case.

```text
{{ toLower "FOOBAR" }}

foobar
```

## `toString` *input*

`toString` returns the string representation of *input*.

```text
{{ toString 10 }}
```

## `toUpper` *string*

`toUpper` returns *string* with all letters converted to upper case.

```text
{{ toUpper "foobar" }}

FOOBAR
```

## `trimSpace` *string*

`trimSpace` returns *string* with all spaces removed.

```text
{{ "    foobar    " | trimSpace }}

foobar
```
