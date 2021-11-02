# logmerge
A Log merging tool for linux.

# How to build
`make build`

# How to run
`--files or -f` will allow you to specify multiple log files (comma-seperated)
`--output or -o` will allow you to specify an output file (default: merge.log)

# Pattern
Currently it matches the following regex: `\d{4}[-\/]\d{2}[-\/]\d{2}[T_ ]*[ ]?\d{1,}:\d{1,}:\d{1,}Z*,*[0-9]*`
you can find this in `cmd/logmerge.go`

## Example:
`bin/logmerge merge --files log1,log2,log3 -o output.log`

##Output
```
INFO[0000] In: [origLogs/colo2.log origLogs/colo3.log origLogs/colo4.log]
INFO[0000] Out: test.log
```

