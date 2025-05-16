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

// initializeNils recursively traverses the value and replaces nil slices and maps with empty ones.
// It respects Go's reflection rules regarding unexported fields:
//   - Exported fields can be read and modified
//   - Unexported fields can be read but not modified (without unsafe)
//   - Unexported pointer fields can be traversed (we can follow the pointer) even though
//     the pointer itself cannot be modified
//   - The fields inside an unexported pointer field cannot be modified either, as they
//     belong to a struct that is not addressable through reflection
func initializeNils(v reflect.Value, visited map[uintptr]bool) {
	if checkVisited(v, visited) {
		return
	}

	// If we somehow received an invalid (zero) reflect.Value, abort early.
	// This can happen when the value originated from an untyped nil stored
	// inside an interface{} or map[*,interface{}].  Calling any method that
	// introspects such a value (Kind, Type, Interface, etc.) would panic, so
	// we must return immediately instead.
	if !v.IsValid() {
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
			if v.CanSet() {
				v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			}
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
			if v.CanSet() {
				v.Set(reflect.MakeMap(v.Type()))
			}
			break
		}

		// Recursively iterate over map items.
		iter := v.MapRange()
		for iter.Next() {
			val := iter.Value()

			// If the value is invalid (untyped nil stored in interface{}), skip.
			if !val.IsValid() {
				continue
			}

			// Map element (value) can't be set directly; we need an addressable copy.
			elemType := val.Type()
			subv := reflect.New(elemType).Elem()

			// Copy its original value.
			subv.Set(val)

			// Replace nil slices and maps inside.
			initializeNils(subv, visited)

			// And set the replacement back in the map.
			v.SetMapIndex(iter.Key(), subv)
		}

	case reflect.Interface:
		// Dereference interface{}.
		if v.IsNil() {
			break
		}

		valueUnderInterface := reflect.ValueOf(v.Interface())
		if !valueUnderInterface.IsValid() {
			return
		}

		elemType := valueUnderInterface.Type()
		subv := reflect.New(elemType).Elem()
		subv.Set(valueUnderInterface)

		initializeNils(subv, visited)

		if v.CanSet() {
			v.Set(subv)
		}

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
			fieldType := v.Type().Field(i)

			if fieldType.IsExported() {
				// Process exported fields normally - these can be both read and modified
				initializeNils(field, visited)
			} else if field.Kind() == reflect.Ptr && !field.IsNil() {
				// Handle unexported pointer fields:
				// Even though the field itself is unexported (and we can't modify the pointer),
				// we can follow the pointer to process the value it points to.
				// However, we can't modify fields inside this dereferenced value
				// because the struct itself is not addressable through reflection.
				initializeNils(field.Elem(), visited)
			}
			// Skip all other unexported fields as we can't modify them without using unsafe
		}
	}
}

// checkVisited tracks values we've already processed to avoid infinite recursion
// in cyclic data structures.
func checkVisited(v reflect.Value, visited map[uintptr]bool) bool {
	if !v.IsValid() {
		return false
	}

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
