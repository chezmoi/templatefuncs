# templatefuncs

Package templatefuncs provides a rich set of
[text/template](https://pkg.go.dev/text/template) functions.

templatefuncs is a modern alternative to
[github.com/masterminds/sprig](https://github.com/masterminds/sprig) with the
following goals:

* Flexible, practical typing of user-equivalent types (e.g. functions that
  accept `string`s also accept `fmt.Stringer`s and `[]byte`s).
* Correct argument order, compatible with text/template's pipelines where the
  most variable argument is passed last (e.g. so you can write `dict "key"
  "value" | hasKey "key"` instead of `hasKey (dict "key" value") "key"`).
* Idiomatic Go naming conventions (e.g. `toJSON`, not `toJson`).
* Structure-preserving transformations (e.g. `toLower` converts a `string` to a
  `string` and also converts a `[]string` to a `[]string`).
* Linkable documentation for individual template functions (so you can direct
  users to the documentation for a single function, not just a page of
  functions).
* Exported documentation which you can include in your own project (so you can
  include a full list of template functions that your project supports).
* Actively maintained.

templatefuncs explicitly is *not* backwards compatible with
github.com/masterminds/sprig.

templatefuncs is currently in the experimental stage and is not suitable for
production use.

## License

MIT
