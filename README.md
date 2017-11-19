# gLISP
Implementation of the programming language Common Lisp written in Go

## Proposal syntax

```
(package
  ; load file on this code
  (load 'lerolero) ; ./lerolero.gl

  ; import package lerolero and used methods
  (import test 'lerolero.test) ; ./lerolero/test.lg
  (print (test.Test))

  ; set vars
  (var a 1)
  (print (a))

  ; create function
  (def hi
    ('Hello))
  (print hi)

  ; create function (by lambda)
  (var area
    (lambda (r) (* 3.141592653 (* r r))))
  (print (area 10.0))
)
```
