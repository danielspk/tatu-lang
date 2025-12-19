# Tatu-Lang

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![CI](https://github.com/danielspk/tatu-lang/workflows/CI/badge.svg)](https://github.com/danielspk/tatu-lang/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/danielspk/tatu-lang)](https://github.com/danielspk/tatu-lang/releases)

![](assets/tatu.png)

> The _Tatu_ mascot was generated with AI 🤖.

---

## What is _Tatu_

_Tatu_ is a toy language developed for general scripting that can be embedded within **Go** applications.

---

## Features

_Tatu_ is a programming language with the following features:

- ✅ Based on _s-expressions_.
- ✅ Everything is an expression _(no statements)_.
- ✅ The last evaluated expression always returns its value.
- ✅ First-class functions with lexical closures.
- ✅ Lambda functions.
- ✅ File inclusion system.
- ✅ Syntactic sugar support.
- ✅ Explicit tail call optimization (TCO).
- ✅ UTF-8 native support.
- ✅ Pure Go implementation.

> _Tatu_ is not a _Lisp_ dialect and does not support metaprogramming.

### TODOs

- ⚠️ Extend the Standard Library _(Math/String/DateTime/Vector/HashMap/IO)_
- ⚠️ Add a Compiler and Virtual Machine

---

## Syntax

_Tatu_ syntax is similar to classic _Lisp_ dialects:

```scheme
(var fib (lambda (n)
  (if (= n 0)
    0
    (if (= n 1)
      1
      (+ (fib (- n 1)) (fib (- n 2)))))))

(print (fib 10))
```

---

## Architecture

_Tatu_ uses a multi-phase pipeline to execute programs:

### Pipeline Interpreted

```
┌─────────────┐
│ Source File │
└──────┬──────┘
       │
       ▼
    Builder ◄──────┐
       │           │
       ▼           │
    Scanner        │
       │           │
       ▼           │
    Parser         │
       │           │
       ▼           │
   Analyzer        │
       │           │
       ▼           │
     Sugar         │
       │           │
       ▼           │
      AST          │
       │          yes
       ▼           │
 (include expr) ───┘
       │
       │ no
       ▼
    Final AST
       │
       ▼
   Interpreter
       │
       ▼
┌─────────────┐
│    Result   │
└─────────────┘
```

### Pipeline Compiled

> TODO

---

## Grammar

The _Tatu_ BNF is extremely simple.

```bnf
<program>      ::= (<expr> | <comment>)*

<expr>         ::= <atom>
                | <list>

<atom>         ::= <number>
                | <string>
                | <boolean>
                | <symbol>
                | "nil"

<list>         ::= "(" <list-body> ")"

<list-body>    ::= <expr>*
                | <primitive>

<primitive>    ::= <include>
                 | <block>
                 | <definition>
                 | <assignment>
                 | <conditional>
                 | <while>
                 | <lambda>
                 | <recur>
                 | <vector>
                 | <hash-map>
                 | <print>

<include>      ::= "include" <string>
<block>        ::= "begin" <expr>+
<definition>   ::= "var" <identifier> <expr>
<assignment>   ::= "set" <identifier> <expr>
<conditional>  ::= "if" <expr> <expr> <expr>?
<while>        ::= "while" <expr> <expr>
<lambda>       ::= "lambda" "(" <identifier>* ")" <expr>
<recur>        ::= "recur" <expr>+
<vector>       ::= "vector" <expr>*
<hash-map>     ::= "map" <key-value>*
<print>        ::= "print" <expr>*

<key-value>    ::= (<identifier> | <string>) <expr>

<comment>      ::= ";" [^\n]*
<number>       ::= ("-")? <digit>+ ("." <digit>+)?
<symbol>       ::= <identifier> | <operator>
<string>       ::= "\"" <character>* "\""
<boolean>      ::= "true" | "false"

<digit>        ::= "0" | "1" | ... | "9"
<letter>       ::= "a" | ... | "z" | "A" | ... | "Z"
<character>    ::= [^"\\] | "\\\\" | "\\\"" | "\\n" | "\\r" | "\\t"
<identifier>   ::= (<letter> | "_") (<letter> | <digit> | "-" | "_" | "?" | ":")*
<operator>     ::= ("+" | "-" | "*" | "/" | "=" | ">" | "<" | "!" | "&" | "|")+

```

> Although the BNF could lack primitives since everything is ultimately an expression, they are included to provide
> the parser with greater control to detect invalid syntax.

---

## Platforms

_Tatu_ is only available for 64-bit architectures. Distributions exist for _Linux_, _Mac_, and _Windows_.

---

## Learning Resources

_Tatu_ is a toy language that emerged as an experiment to put programming language design and development concepts into
practice for educational purposes. Below are some resources of interest:

- **Courses:**
  - [Compilers, Interpreters & Formal Languages - Pikuma](https://pikuma.com/courses/create-a-programming-language-compiler)
  - [Programming Languages Ultimate Bundle - Dmitry Soshnikov](https://www.dmitrysoshnikov.education/p/programming-languages-ultimate-bundle-3rd-edition)
- **Books:**
  - [Writing An Interpreter In Go - Thorsten Ball](https://interpreterbook.com/)
  - [Crafting Interpreters - Robert Nystrom](https://craftinginterpreters.com/)
  - [Build Your Own Lisp - Daniel Holden](https://www.buildyourownlisp.com/)
  - [The Go Compiler Builder's Handbook: From Lexical Analysis to Code Generation - Jayden Reed](https://www.amazon.com/-/es/Jayden-Reed/dp/B0DJM2C65G)
- **Articles:**
  - [How to Write a Lisp Interpreter in Python - Peter Norvig](https://norvig.com/lispy.html)
  - [An Even Better Lisp Interpreter in Python - Peter Norvig](https://www.norvig.com/lispy2.html)
  - [Lisp in 99 lines of C and how to write one yourself - Robert van Engelen](https://github.com/Robert-van-Engelen/tinylisp/blob/main/tinylisp.pdf)
  - [How to Write a Lisp Interpreter in JavaScript - Chidi Williams](https://www.chidiwilliams.com/posts/how-to-write-a-lisp-interpreter-in-javascript)
  - [Let's Build A Simple Interpreter series - Ruslan Spivak](https://ruslanspivak.com/lsbasi-part1/)
  - [Mal (make-a-lisp) - Joel Martin](https://github.com/kanaka/mal)
- **Resources:**
  - [Awesome compilers](https://aalhour.com/awesome-compilers/)

---

## License

This software is distributed under the _MIT_ license. See the [LICENSE](LICENSE) file.
