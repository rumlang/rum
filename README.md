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

* To re-generate `y.go` from the grammar, assuming that you are in the `glop/` directory:
```bash
go tool yacc -o parser/y.go parser/glop.y
```

* To run tests, for example for the `runner` module:
```bash
go test github.com/palats/glop/runner
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
