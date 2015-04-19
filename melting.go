package melting

import (
	"errors"
	"fmt"
	"reflect"
)

// Melt assigns a source value to a destination.
// If source and destination are not structs, they must
// have the same type. If not, melting is applied.
func Melt(src, dest interface{}) error {
	srcEl := reflect.ValueOf(src).Elem()
	destEl := reflect.ValueOf(dest).Elem()
	if destEl.Kind() == reflect.Struct {
		srcType := srcEl.Type()
		for i := 0; i < srcEl.NumField(); i++ {
			fieldName := srcType.Field(i).Name
			if destField := destEl.FieldByName(fieldName); destField.IsValid() {
				srcField := srcEl.Field(i)
				err := meltValue(srcField, destField)
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return meltValue(srcEl, destEl)
	}
}

func meltValue(src, dest reflect.Value) error {
	if !dest.CanSet() {
		return errors.New(fmt.Sprintf("destination field %v is not assignable", dest))
	}
	if !dest.Type().AssignableTo(src.Type()) {
		return errors.New(fmt.Sprintf("cannot assign type %v to %v", src.Type(), dest.Type()))
	}
	fmt.Printf("meltValue from %v to %v\n", src, dest)
	dest.Set(src)
	return nil
}
