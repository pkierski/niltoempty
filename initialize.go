package niltoempty

import (
	"reflect"
)

// Initialize traverses any addressable entity and replaces all nil maps and slices
// with empty map and slices respectively.
//
// Because input object have to be addressable in order to make changes Initialize
// panics when non-adresable object is passed as argument.
//
// Because pointer to element is usually used for modeling optional fields
// nil pointers to the map or slices are left untouched.
func Initialize(obj interface{}) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		panic("niltoempty: expected pointer")
	}

	initializeNils(v)

	return obj
}

func initializeNils(v reflect.Value) {
	// Dereference pointer(s).
	for (v.Kind() == reflect.Ptr) && !v.IsNil() {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice:
		// Initialize a nil slice.
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			return
		}

		// Recursively iterate over slice items.
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			initializeNils(item)
		}

	case reflect.Map:
		// Initialize a nil map.
		if v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
			return
		}

		// Recursively iterate over map items.
		iter := v.MapRange()
		for iter.Next() {
			// Map elements (values) aren't addressable.

			// we have to alloc addressable replacement for it
			elemType := iter.Value().Type()
			subv := reflect.New(elemType).Elem()
			// copy its original value
			subv.Set(iter.Value())

			// replace nil slices and maps inside
			initializeNils(subv)

			// and set the replacement back in map
			v.SetMapIndex(iter.Key(), subv)
		}

	case reflect.Interface:
		// Dereference interface{}.
		if v.IsNil() {
			return
		}
		valueUnderInterface := reflect.ValueOf(v.Interface())
		elemType := valueUnderInterface.Type()
		subv := reflect.New(elemType).Elem()
		subv.Set(valueUnderInterface)

		initializeNils(subv)

		v.Set(subv)

	// Recursively iterate over array elements.
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			initializeNils(elem)
		}

	// Recursively iterate over struct fields.
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			initializeNils(field)
		}

	}
}
