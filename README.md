A lisp interpreter in Go
------------------------
This is mostly a toy project and not aiming to be a full lisp interpreter in Go. The base is largely inspired by [Lispy by Peter Norvig](http://norvig.com/lispy.html) - though it does not try to be minimalistic.

Building & Running
------------------
This uses standard [Go project hierarchy](see http://golang.org/doc/code.html). It is supposed to go in `$GOPATH/src/github.com/palats/glop`.

* You will need to following libraries:
```bash
go get github.com/GeertJohan/go.linenoise
go get github.com/golang/glog
```

* To re-generate `y.go` from the grammar, assuming that you are in the `glop/` directory:
```bash
go tool yacc -o parser/y.go parser/glop.y
```

* To run tests, for example for the `runner` module:
```bash
go test github.com/palats/glop/runner
```
