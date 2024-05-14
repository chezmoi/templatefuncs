# String Functions

## `contains` *substring* *string*

`contains` returns whether *substring* is in *string*.

```text
{{ "abc" | contains "ab" }}

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

## `toLower` *string*

`toLower` returns *string* with all letters converted to lower case.

```text
{{ toLower "FOOBAR" }}

foobar
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
