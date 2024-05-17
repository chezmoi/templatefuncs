# Boolean Functions

## `eqFold` *string1* *string2* [*extraStrings*...]

`eqFold` returns the boolean truth of comparing *string1* with *string2*
and any number of *extraStrings* under Unicode case-folding.

```text
{{ eqFold "föö" "FOO" }}

true
```
