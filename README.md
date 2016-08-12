# gopprof
A go pprof wrapper which provided full readline user experience

### Get
```
$ go get github.com/chzyer/gopprof
$ gopprof [options] [binary] <profile source> ...
# just like anything with "go tool pprof"
```

[![asciicast](https://asciinema.org/a/c2ijwhevuxhic8j1n1eylgqx3.png)](https://asciinema.org/a/c2ijwhevuxhic8j1n1eylgqx3)

```
$ gopprof http://localhost:6060/debug/pprof/heap
Fetching profile from http://localhost:6060/debug/pprof/heap
Saved profile in /Users/xxx/pprof/pprof.localhost:6060.inuse_objects.inuse_space.028.pb.gz
Entering interactive mode (type "help" for commands)
(gopprof) _
```
