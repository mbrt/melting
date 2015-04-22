// Copyright 2015 Michele Bertasi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
	Package melting provides an utility to merge structures of differennt
	types. Fields of the source structure are assigned to fields of the
	destination structure matching by field names.

	The destination parameter must be a pointer to a structure, because
	its fields will be overridden by fields of the source structure.

	Given a field F: if F is present in the source and destination structures,
	the source value will override the destination value; if F is present in
	the source structure but not in destination, the field will be ignored;
	if F is present in the destination structure, but not in the source,
	the destination will preserve its value. Embedded sub-structures are
	supported, and the same algorithm is applied for them.

	For example:
		type Source struct {
			int    F1
			string F2
		}
		type Dest struct {
			int    F1
			string F2
			real   F3
		}

		func Example() {
			s := Source{F1: 3, F2: "a"}
			d := Dest{F1: 4, F2: "b", F3: 3.0}
			melting.Melt(s, &d)
		}

	After the Melt call, source will not change, while
	destination will assume this value:

		Dest{F1: 3, F2: "a", F3: 3.0}
*/
package melting

import (
	"errors"
	"fmt"
	"reflect"
)

// Melt assigns a source value to a destination.
// If not, melting is applied. The fields of the destination
// struct will get the value of the source fields, for those
// they have in common. If those fields have different types,
// an error will be returned.
func Melt(src, dest interface{}) error {
	// check dest ptr
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("dest value %v is not Ptr", dest))
	}
	destEl := reflect.ValueOf(dest).Elem()

	// handle optional src ptr
	srcEl := reflect.ValueOf(src)
	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		srcEl = srcEl.Elem()
	}

	return meltValue(srcEl, destEl)
}

// Filter is an interface that can be used to instrument the melting
// function to not consider certain fields for the melt.
type Filter interface {
	// ConsiderField returns true if the given field name have
	// to be considered for the melting. The source and destination
	// values are also provided.
	ConsiderField(name string, src, dest reflect.Value) bool
}

// MeltWithFilter is just like Melt, but with a user supplied filter,
// that allows to ignore certain fields in the melt.
func MeltWithFilter(src, dest interface{}, filter Filter) error {
	return nil
}

func meltValue(src, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.Struct:
		return meltStruct(src, dest)
	default:
		return meltAssignable(src, dest)
	}
}

func meltStruct(src, dest reflect.Value) error {
	srcType := src.Type()
	for i := 0; i < src.NumField(); i++ {
		fieldName := srcType.Field(i).Name
		if destField := dest.FieldByName(fieldName); destField.IsValid() {
			srcField := src.Field(i)
			err := meltValue(srcField, destField)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func meltAssignable(src, dest reflect.Value) error {
	if !dest.CanSet() {
		return errors.New(fmt.Sprintf("destination field %v is not assignable", dest))
	}
	if !dest.Type().AssignableTo(src.Type()) {
		return errors.New(fmt.Sprintf("cannot assign type %v to %v", src.Type(), dest.Type()))
	}
	dest.Set(src)
	return nil
}
