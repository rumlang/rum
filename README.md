About Glop
----------
This project implements a simple lisp interpreter in Go.

It does not try to provide a full lisp/scheme implementation. This is mostly a toy project; the only focus so far was to try to have some error detection and reporting, something which is often overlooked in parsers&interpreters tutorials.

It tries to map as closely as possible to Go, relying on it when possible for most mechanisms.

The initial code was inspired by [Lispy by Peter Norvig](http://norvig.com/lispy.html) - though the actual code has nothing in common and does not try to be minimalistic.

Building & Running
------------------
This uses standard [Go project hierarchy](see http://golang.org/doc/code.html). It is supposed to be placed in `$GOPATH/src/github.com/palats/glop`.

* You will need to following libraries:
```bash
go get github.com/GeertJohan/go.linenoise
go get github.com/golang/glog
```

* Then, you can run the REPL:
```
go run $GOPATH/src/github.com/palats/glop/glop.go
```

* Or you can execute a glop file:
```
go run $GOPATH/src/github.com/palats/glop/glop.go $GOPATH/src/github.com/palats/glop/example.glop
```

* To run all tests, from the glop/ directory:
```bash
go test ./...
```

Features
--------

File `example.glop` shows a few examples.

Available list of functions is in `runtime/runtime.go`, in the `NewContext` function. Currently:

* Basic lisp: `begin`, `quote`, `define`, `set!`, `if`, `lambda`, `cons`, `car`, `cdr`

* Constants: `true`, `false`

* Integers and float numbers are supported, but scientific notation is not.

* Operators: `==`, `!=`, `<`, `<=`, `>`, `>=`, `+`, `-`, `*` ; they look at the
  first argument to see whether they are operating on integer or floats.

* Misc: `print`, `panic` (triggers a panic)


Example:
```
In [0]: (define area (lambda (r) (* 3.141592653 (* r r))))
Out [0]: <runtime.Internal>(runtime.Internal)(0x457d00)
In [1]: (area 3.0)
Out [1]: <float64>28.274333877
In [2]: (define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))
Out [2]: <runtime.Internal>(runtime.Internal)(0x457d00)
In [3]: (fact 10)
Out [3]: <int64>3628800
```


Project license
---------------
This project is under the Apache License version 2.0.

If not explicitely mentionned otherwise, all files in this project are under those terms:

    Copyright 2014 Pierre Palatin

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
