# List Functions

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

## `prefixLines` *prefix* *list*

`prefixLines` returns a string consisting of each item in *list*
with the prefix *prefix*.

```text
{{ list "this is" "a multi-line" "comment" | prefixLines "# " }}

# this is
# a multi-line
# comment
```
