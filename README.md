 melting
=========================================================

[![GoDoc](https://godoc.org/github.com/mbrt/melting?status.png)](https://godoc.org/github.com/mbrt/melting) [![Build Status](https://travis-ci.org/mbrt/melting.svg)](https://travis-ci.org/mbrt/melting) [![Coverage Status](https://coveralls.io/repos/mbrt/melting/badge.svg)](https://coveralls.io/r/mbrt/melting) 

Go package allowing to melt two heterogeneous structures. This is useful for merge configuration structures coming from different sources (e.g. command line vs. configuration file). You can override the configuration file options with the command line flags, even if the types (and so the number and positioning of the fields) are different.

For example:

```go
type Source struct {
    F1  int
    F2  string
}
type Dest struct {
    F1  int
    F3  float32
    F2  string
}

s := Source{F1: 3, F2: "a"}
d := Dest{F1: 4, F2: "b", F3: 3.1}

// Source will override Dest for fields 'F1', 'F2';
// field 'F3' will be ignored.
// Note that fields order does not count.
err := Melt(s, &d)
if err != nil {
    log.Fatalf("cannot assign source to dest, error %v", err)
}

fmt.Printf("Source%v\nDest%v\n", s, d)
```

will output:

```
Source{3 a}
Dest{3 3.1 a}
```
