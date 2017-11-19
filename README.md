# gLISP
Free software environment for statistical computing

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

  ; create loops with do
  (do (print 'hello)) ; will loop while expression or funcrion return false
  (do area (10.0 20.0 30.0) ; will interact on the list elements
  (do (var a 1)
    (= a 10)
      (var a (+ a 1)
      (print a)))

  ; create if 
  (if (= a 10)        ; if a is equal to 10 
    (print 'Hello))   ; print Hello

  (if (= a 'good)     ; if a is equal to "good"
    (print 'good)     ; print "good"
    (else             ; otherwise
      (print 'bad)))  ; print "bad"

  (if (= a 'good)     ; if a is equal to "good"
    (print 'good)     ; print "good"
    (else (= a 'bad)  ; if not and is equal to "bad"
      (print 'bad))   ; print "bad"
    (else             ; otherwise
      (print 'ugly))) ; print "ugly"

```
