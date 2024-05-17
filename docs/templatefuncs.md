# Template Functions

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

## `stat` *path*

`stat` returns a map representation of executing
[`os.Stat`](https://pkg.go.dev/os#Stat) on *path*.

```text
{{ (stat "some/file").type }}

file
```
