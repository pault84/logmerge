# logmerge
A Log merging tool for linux.

# How to build
`make build`

# Pattern
Currently it matches the following regex: `\d{4}[-\/]\d{2}[-\/]\d{2}[T_ ]*[ ]?\d{1,}:\d{1,}:\d{1,}Z*,*[0-9]*`
you can find this in `cmd/logmerge.go`

## Example:
`bin/logmerge log1 log2 log3 > output.log`

## Output
```
INFO[0000] In: [origLogs/colo2.log origLogs/colo3.log origLogs/colo4.log]
```

