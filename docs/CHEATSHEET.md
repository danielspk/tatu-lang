# Tatu Language Cheatsheet

## Basic Syntax

### Atoms
```lisp
42          ; number
3.14        ; number
"hello"     ; string
true        ; boolean
false       ; boolean
nil         ; nil
symbol      ; symbol/identifier
```

### Lists (S-expressions)
```lisp
(operator arg1 arg2 ...)

(+ 1 2 3)

(print "hello")
```

### Comments
```lisp
; This is a comment
```

## Variables

```lisp
(var name value)           ; define variable

(set name new-value)       ; assign to variable
```

## Functions

```lisp
(lambda (params) body)     ; anonymous function

(def name (params) body)   ; define function (sugar)
```

## Control Flow

```lisp
(if condition 
    then 
    [else])

(while condition
    body)

(for init condition increment
    body)

(switch
    (cond1 result1)
    (cond2 result2)
    (default result3))

(begin 
    expr1
    expr2 ...)

(recur arg1 arg2 ...)   ; tail recursion
```

## Operators

```lisp
; Arithmetic
(+ a b ...)            ; addition
(- a b ...)            ; subtraction
(* a b ...)            ; multiplication
(/ a b ...)            ; division
(% a b)                ; modulo
(- x)                  ; unary negation

; Comparison
(= a b)
(< a b)
(<= a b)
(> a b)
(>= a b)

; Logical
(and a b ...)          ; logical and
(or a b ...)           ; logical or
(not a)                ; logical negation
```

## Data Structures

### Vectors
```lisp
(vector 1 2 3)                     ; (1 2 3)
(vector)                           ; ()

(vec:get (vector 10 20 30) 0)      ; 10
(vec:set (vector 1 2 3) 1 99)      ; (1 99 3)
```

### Maps
```lisp
(map "name" "John" "age" 30)       ; {name: "John", age: 30}
(map)                              ; {}

(map:get (map "x" 10 "y" 20) "x")  ; 10
(map:set m "key" "value")          ; {key: "value"}
```

## Sugar Syntax

```lisp
; function definition:
(def name (params)
    body)

; expands to:
(var name (lambda (params)
    body))
```

```lisp
; for loop:
(for init cond inc
    body)

; expands to:
(begin
    init 
    (while cond 
        (begin
            body
            inc)))
```

```lisp
; switch:
(switch
    (cond1 result1)
    (cond2 result2)
    (default result3))

; expands to
(if cond1 result1
(if cond2 result2
          result3))

```

## Module System

```lisp
(include "path/to/file.tatu")
```

## Core Builtins

### I/O
```lisp
(print x ...)
```

### Type Checking
```lisp
(is-bool x)
(is-number x)
(is-int x)
(is-string x)
(is-vector x)
(is-map x)
(is-nil x)
(is-function x)
```

### Type Conversion
```lisp
(to-string x)
(to-number x)
(to-bool x)
```

## Standard Library

### Math

| Function | Description |
|----------|-------------|
| `(math:pi)` | π constant |
| `(math:e)` | e constant |
| `(math:abs x)` | Absolute value |
| `(math:floor x)` | Floor |
| `(math:ceil x)` | Ceiling |
| `(math:round x)` | Round |
| `(math:sqrt x)` | Square root |
| `(math:pow x y)` | Power |
| `(math:sin x)` | Sine |
| `(math:cos x)` | Cosine |
| `(math:tan x)` | Tangent |
| `(math:log x)` | Natural logarithm |
| `(math:exp x)` | e^x |
| `(math:min x y)` | Minimum |
| `(math:max x y)` | Maximum |
| `(math:between x min max)` | Check if x in range [min, max] |
| `(math:rand min max)` | Random integer in range |

### String

| Function | Description |
|----------|-------------|
| `(str:len s)` | Length |
| `(str:concat s1 s2 ...)` | Concatenate |
| `(str:split s sep)` | Split by separator |
| `(str:join vec sep)` | Join with separator |
| `(str:slice s start end)` | Substring |
| `(str:contains s substr)` | Check contains |
| `(str:starts s prefix)` | Check starts with |
| `(str:ends s suffix)` | Check ends with |
| `(str:index s substr)` | Find index |
| `(str:upper s)` | Uppercase |
| `(str:lower s)` | Lowercase |
| `(str:trim s)` | Trim whitespace |
| `(str:replace s old new)` | Replace all |
| `(str:repeat s n)` | Repeat n times |
| `(str:reverse s)` | Reverse |

### Vector

| Function | Description |
|----------|-------------|
| `(vec:len v)` | Length |
| `(vec:get v i)` | Get element at index |
| `(vec:set v i val)` | Set element at index |
| `(vec:push v val)` | Append element |
| `(vec:pop v)` | Remove last element |
| `(vec:concat v1 v2)` | Concatenate |
| `(vec:slice v start end)` | Subvector |
| `(vec:find v val)` | Find index of value |
| `(vec:contains v val)` | Check contains |
| `(vec:delete v i)` | Delete at index |
| `(vec:reverse v)` | Reverse |
| `(vec:sort v)` | Sort ascending |

### Map

| Function | Description |
|----------|-------------|
| `(map:len m)` | Number of keys |
| `(map:get m key)` | Get value |
| `(map:get-in m path)` | Deep access with path vector |
| `(map:set m key val)` | Set key-value |
| `(map:has m key)` | Check key exists |
| `(map:delete m key)` | Delete key |
| `(map:keys m)` | Get all keys |
| `(map:values m)` | Get all values |
| `(map:merge m1 m2)` | Merge maps |

### Time

| Function | Description |
|----------|-------------|
| `(time:now)` | Current time |
| `(time:unix t)` | Unix timestamp |
| `(time:year t)` | Get year |
| `(time:month t)` | Get month |
| `(time:day t)` | Get day |
| `(time:hour t)` | Get hour |
| `(time:minute t)` | Get minute |
| `(time:second t)` | Get second |
| `(time:format t layout)` | Format time |
| `(time:parse layout s)` | Parse time |
| `(time:add t duration)` | Add duration |
| `(time:sub t duration)` | Subtract duration |
| `(time:diff t1 t2)` | Difference |
| `(time:is-leap year)` | Check leap year |

### JSON

| Function | Description |
|----------|-------------|
| `(json:encode val)` | Encode to JSON |
| `(json:decode s)` | Decode from JSON |

### File System

| Function | Description |
|----------|-------------|
| `(fs:read path)` | Read file |
| `(fs:write path content)` | Write file |
| `(fs:append path content)` | Append to file |
| `(fs:delete path)` | Delete file |
| `(fs:exists path)` | Check exists |
| `(fs:list path)` | List directory |
| `(fs:mkdir path)` | Create directory |
| `(fs:move src dst)` | Move/rename |
| `(fs:is-dir path)` | Check if directory |
| `(fs:basename path)` | Get basename |
| `(fs:size path)` | Get file size |
| `(fs:temp-dir)` | Get temp directory |
| `(fs:read-lines path)` | Read lines |

### Regex

| Function | Description |
|----------|-------------|
| `(regex:matches s pattern)` | Check if matches |
| `(regex:find s pattern)` | Find first match |
| `(regex:replace s pattern repl)` | Replace all |

## Examples

```lisp
; Factorial with tail recursion

(def factorial (n acc)
  (if (= n 0)
    acc
    (recur (- n 1) (* acc n))))

(factorial 5 1)
```

```lisp
; Filter numbers in range

(def filter-range (numbers min max)
  (var result (vector))
  (var i 0)

  (while (< i (vec:len numbers))
    (begin
      (var num (vec:get numbers i))
      (if (math:between num min max)
        (vec:push result num))
      (set i (+ i 1))))
  result)

(filter-range (vector 1 5 10 15 3 20) 5 15)
```

---

> *This cheatsheet was AI-generated and may contain errors. Please refer to official documentation for accurate information.*
