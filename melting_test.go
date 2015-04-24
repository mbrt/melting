// Copyright 2015 Michele Bertasi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package melting

import (
	"reflect"
	"testing"
)

func TestField(t *testing.T) {
	src := 3.14
	dest := 15.0
	err := Melt(&src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src != 3.14 {
		t.Fatalf("expected src: %v, got %v", 3.14, src)
	}
	if dest != 3.14 {
		t.Fatalf("expected dest: %v, got %v", 3.14, dest)
	}
}

func TestErrField(t *testing.T) {
	src := 3.14
	dest := 15
	err := Melt(&src, &dest)
	if err == nil {
		t.Fatalf("different type assignment is possible")
	}
}

func TestErrNotRefField(t *testing.T) {
	src := 3.14
	dest := 15.0
	err := Melt(&src, dest)
	if err == nil {
		t.Fatalf("No error reported with dest not pointer")
	}
}

func TestNotRefSrcField(t *testing.T) {
	src := 3.14
	dest := 15.0
	err := Melt(src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src != 3.14 {
		t.Fatalf("expected src: %v, got %v", 3.14, src)
	}
	if dest != 3.14 {
		t.Fatalf("expected dest: %v, got %v", 3.14, dest)
	}
}

type Simple struct {
	F1 string
	F2 bool
	F3 int
}

func TestSameStruct(t *testing.T) {
	src := Simple{F1: "a", F2: true, F3: 7}
	dest := Simple{F1: "b", F2: false, F3: 8}

	err := Melt(&src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src.F1 != "a" || src.F2 != true || src.F3 != 7 {
		t.Fatalf("changed source struct in %v", src)
	}
	if dest.F1 != "a" || dest.F2 != true || dest.F3 != 7 {
		t.Fatalf("expected dest struct: %v, got: %v", src, dest)
	}
}

func TestNotRefSrcStruct(t *testing.T) {
	src := Simple{F1: "a", F2: true, F3: 7}
	dest := Simple{F1: "b", F2: false, F3: 8}

	err := Melt(src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src.F1 != "a" || src.F2 != true || src.F3 != 7 {
		t.Fatalf("changed source struct in %v", src)
	}
	if dest.F1 != "a" || dest.F2 != true || dest.F3 != 7 {
		t.Fatalf("expected dest struct: %v, got: %v", src, dest)
	}
}

type Bigger struct {
	F1 string
	F2 bool
	F3 int
	F4 uint
}

func TestBiggerSrc(t *testing.T) {
	src := Bigger{F1: "a", F2: true, F3: 7, F4: 8}
	dest := Simple{F1: "b", F2: false, F3: 8}

	err := Melt(&src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src.F1 != "a" || src.F2 != true || src.F3 != 7 {
		t.Fatalf("changed source struct in %v", src)
	}
	if dest.F1 != "a" || dest.F2 != true || dest.F3 != 7 {
		t.Fatalf("expected dest struct: %v, got: %v", src, dest)
	}
}

type RSimple struct {
	F2 bool
	F1 string
	F3 int
}

func TestReorderedStruct(t *testing.T) {
	src := RSimple{F1: "a", F2: true, F3: 7}
	dest := Simple{F1: "b", F2: false, F3: 8}

	err := Melt(&src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src.F1 != "a" || src.F2 != true || src.F3 != 7 {
		t.Fatalf("changed source struct in %v", src)
	}
	if dest.F1 != "a" || dest.F2 != true || dest.F3 != 7 {
		t.Fatalf("expected dest struct: %v, got: %v", src, dest)
	}
}

type Smaller struct {
	F1 string
	F3 int
}

func TestSmallerStruct(t *testing.T) {
	src := Smaller{F1: "a", F3: 7}
	dest := Simple{F1: "b", F2: false, F3: 8}

	err := Melt(&src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src.F1 != "a" || src.F3 != 7 {
		t.Fatalf("changed source struct in %v", src)
	}
	if dest.F1 != "a" || dest.F2 != false || dest.F3 != 7 {
		t.Fatalf("expected dest struct: %v, got: %v", src, dest)
	}
}

type Embedded struct {
	F1 int
	F2 Simple
}

func TestEmbeddedStruct(t *testing.T) {
	src := Embedded{F1: 1, F2: Simple{F1: "a", F2: true, F3: 7}}
	dest := Embedded{F1: 2, F2: Simple{F1: "a", F2: false, F3: 7}}
	err := Melt(&src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src.F1 != 1 || src.F2.F1 != "a" || src.F2.F2 != true || src.F2.F3 != 7 {
		t.Fatalf("changed source struct in %v", src)
	}
	if dest.F1 != 1 || dest.F2.F1 != "a" || dest.F2.F2 != true || dest.F2.F3 != 7 {
		t.Fatalf("expected dest struct: %v, got: %v", src, dest)
	}
}

type EmbeddedBigger struct {
	F1 int
	F2 Bigger
}

func TestEmbeddedBiggerStruct(t *testing.T) {
	src := EmbeddedBigger{F1: 1, F2: Bigger{F1: "a", F2: true, F3: 7, F4: 8}}
	dest := Embedded{F1: 2, F2: Simple{F1: "a", F2: false, F3: 7}}
	err := Melt(&src, &dest)
	if err != nil {
		t.Fatalf("cannot set %v to %v, with error: %v", src, dest, err)
	}
	if src.F1 != 1 || src.F2.F1 != "a" || src.F2.F2 != true || src.F2.F3 != 7 {
		t.Fatalf("changed source struct in %v", src)
	}
	if dest.F1 != 1 || dest.F2.F1 != "a" || dest.F2.F2 != true || dest.F2.F3 != 7 {
		t.Fatalf("expected dest struct: %v, got: %v", src, dest)
	}
}

type filter struct {
	exclude string
}

func (f *filter) Filter(srcField, destField reflect.StructField, src, dest reflect.Value) bool {
	return f.exclude != srcField.Name
}

func NewFilter(exclusion string) *filter {
	return &filter{exclude: exclusion}
}

func TestFilter(t *testing.T) {
	src := Simple{F1: "a", F2: true, F3: 7}
	dest := Simple{F1: "b", F2: false, F3: 8}

	err := MeltWithFilter(&src, &dest, NewFilter("F2"))
	if err != nil {
		t.Fatalf("cannot set %v to %v", src, dest)
	}
	if src.F1 != "a" || src.F2 != true || src.F3 != 7 {
		t.Fatalf("changed source struct in %v", src)
	}
	if dest.F1 != "a" || dest.F2 != false || dest.F3 != 7 {
		t.Fatalf("expected dest struct: %v, got: %v", src, dest)
	}
}
