package niltoempty

import (
	"reflect"
)

// Initialize traverses any addressable entity and replaces all nil maps and slices
// with empty map and slices respectively.
//
// Because input object have to be addressable in order to make changes Initialize
// panics when non-addressable object is passed as argument.
//
// Because pointer to element is usually used for modeling optional fields
// nil pointers to the map or slices are left untouched.
func Initialize(obj interface{}) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		panic("niltoempty: expected pointer")
	}

	initializeNils(v, map[uintptr]bool{})

	return obj
}

func initializeNils(v reflect.Value, visited map[uintptr]bool) {
	if checkVisited(v, visited) {
		return
	}

	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			initializeNils(v.Elem(), visited)
		}
	case reflect.Slice:
		// Initialize a nil slice.
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			break
		}

		// Recursively iterate over slice items.
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			initializeNils(item, visited)
		}

	case reflect.Map:
		// Initialize a nil map.
		if v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
			break
		}

		// Recursively iterate over map items.
		iter := v.MapRange()
		for iter.Next() {
			// Map element (value) can't be set directly.
			// we have to alloc addressable replacement for it
			elemType := iter.Value().Type()
			subv := reflect.New(elemType).Elem()

			// copy its original value
			subv.Set(iter.Value())

			// replace nil slices and maps inside
			initializeNils(subv, visited)

			// and set the replacement back in map
			v.SetMapIndex(iter.Key(), subv)
		}

	case reflect.Interface:
		// Dereference interface{}.
		if v.IsNil() {
			break
		}

		valueUnderInterface := reflect.ValueOf(v.Interface())
		elemType := valueUnderInterface.Type()
		subv := reflect.New(elemType).Elem()
		subv.Set(valueUnderInterface)

		initializeNils(subv, visited)

		v.Set(subv)

	// Recursively iterate over array elements.
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			initializeNils(elem, visited)
		}

	// Recursively iterate over struct fields.
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			initializeNils(field, visited)
		}
	}

}

func checkVisited(v reflect.Value, visited map[uintptr]bool) bool {
	kind := v.Kind()
	if kind == reflect.Map || kind == reflect.Ptr || kind == reflect.Slice {
		if v.IsNil() {
			return false
		}
		p := v.Pointer()
		wasVisited := visited[p]
		visited[p] = true
		return wasVisited
	}
	return false
}
