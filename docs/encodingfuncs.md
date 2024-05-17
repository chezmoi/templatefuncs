# Encoding Functions

## `fromJSON` *jsontext*

`fromJSON` parses *jsontext* as JSON and returns the parsed value.

```text
{{ `{ "foo": "bar" }` | fromJSON }}
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

## `toJSON` *input*

`toJSON` returns a JSON string representation of *input*.

```text
{{ list "foo" "bar" "baz" }}

["foo","bar","baz"]
```
