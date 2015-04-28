package melting

import (
	"fmt"
	"log"
	"reflect"
)

func Example() {
	type Source struct {
		F1 int
		F2 string
	}
	type Dest struct {
		F1 int
		F3 float32
		F2 string
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

	// Output:
	// Source{3 a}
	// Dest{3 3.1 a}
}

type nameFilter struct {
	exclude string
}

func (f *nameFilter) Filter(srcField, destField reflect.StructField, src, dest reflect.Value) bool {
	return f.exclude != srcField.Name
}

func ExampleMeltWithFilter() {
	type Source struct {
		F1 int
		F2 string
	}
	type Dest struct {
		F1 int
		F2 string
		F3 float32
	}

	s := Source{F1: 3, F2: "a"}
	d := Dest{F1: 4, F2: "b", F3: 3.1}

	// Source will override Dest for field 'F1' only;
	// field 'F3' will be ignored because Source does not have it,
	// while 'F2' will be ignored because of the filter.
	//
	// nameFilter is a struct that contains a field name.
	// Its 'Filter' function returns true only if the given field
	// name is different than the stored one.
	err := MeltWithFilter(s, &d, &nameFilter{"F2"})
	if err != nil {
		log.Fatalf("cannot assign source to dest, error %v", err)
	}

	fmt.Printf("Source%v\nDest%v\n", s, d)

	// Output:
	// Source{3 a}
	// Dest{3 b 3.1}
}
