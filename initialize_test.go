package niltoempty_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pkierski/niltoempty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	T struct {
		M map[string]any `json:"m"`
		S []any          `json:"s"`
	}
	TP struct {
		PM *map[string]any `json:"pm"`
		PS *[]any          `json:"ps"`
	}
)

func Example() {
	type T struct {
		M  map[string]any  `json:"m"`
		S  []any           `json:"s"`
		PM *map[string]any `json:"pm"`
		PS *[]any          `json:"ps"`
	}

	var v T

	m1, _ := json.MarshalIndent(v, "", "    ")
	fmt.Println(string(m1))
	m2, _ := json.MarshalIndent(niltoempty.Initialize(&v), "", "    ")
	fmt.Println(string(m2))
	// output
	// {
	//     "m": null,
	//     "s": null,
	//     "pm": null,
	//     "ps": null
	// }
	// {
	//     "m": {},
	//     "s": [],
	//     "pm": null,
	//     "ps": null
	// }
}

func TestNonPointer(t *testing.T) {
	assert.Panics(t, func() {
		var v T
		niltoempty.Initialize(v)
	}, "struct as value")

	assert.Panics(t, func() {
		var v string
		niltoempty.Initialize(v)
	}, "string as value")

	assert.Panics(t, func() {
		var v map[string]any
		niltoempty.Initialize(v)
	}, "map as value")

	assert.Panics(t, func() {
		var v []any
		niltoempty.Initialize(v)
	}, "slice as value")
}

func TestSlice(t *testing.T) {
	t.Run("as root", func(t *testing.T) {
		v := []any(nil)

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `[]`, string(b))
	})

	t.Run("in struct", func(t *testing.T) {
		var v T

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"m":{},"s":[]}`, string(b))
	})

	t.Run("in slice", func(t *testing.T) {
		v := make([][]any, 3)

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `[[],[],[]]`, string(b))
	})

	t.Run("in map", func(t *testing.T) {
		v := map[string][]any{
			"a": nil,
			"b": nil,
		}

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"a":[],"b":[]}`, string(b))
	})

	t.Run("in array", func(t *testing.T) {
		v := [2][]any{}

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `[[],[]]`, string(b))
	})
}

func TestMap(t *testing.T) {
	t.Run("as root", func(t *testing.T) {
		var v map[string]string

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{}`, string(b))
	})

	t.Run("in struct", func(t *testing.T) {
		var v T
		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"m":{},"s":[]}`, string(b))
	})

	t.Run("in slice", func(t *testing.T) {
		v := make([]map[string]any, 3)

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `[{},{},{}]`, string(b))
	})

	t.Run("in map", func(t *testing.T) {
		v := map[string]map[string]any{
			"a": nil,
			"b": nil,
		}

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"a":{},"b":{}}`, string(b))
	})

	t.Run("in array", func(t *testing.T) {
		v := [2]map[string]any{}

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `[{},{}]`, string(b))
	})

	t.Run("as interface{}", func(t *testing.T) {
		var (
			v any
			m map[string]any
		)
		v = any(m)

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{}`, string(b))
	})

	t.Run("in struct as inteface{}", func(t *testing.T) {
		type S struct {
			S any `json:"s"`
			E any `json:"e"`
		}
		var m map[string]string
		v := map[string]any{
			"a": any(S{S: any(m)}),
		}

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"a":{"s":{},"e":null}}`, string(b))
	})
}

func TestPointers(t *testing.T) {
	t.Run("leave nil pointers", func(t *testing.T) {
		var v TP
		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"pm":null,"ps":null}`, string(b))
	})

	t.Run("update pointed", func(t *testing.T) {
		var (
			emptySlice []any
			emptyMap   map[string]any
		)
		v := TP{
			PM: &emptyMap,
			PS: &emptySlice,
		}
		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"pm":{},"ps":[]}`, string(b))
	})
}

func TestStruct(t *testing.T) {
	t.Run("as root", func(t *testing.T) {
		type (
			Inner struct {
				S []string `json:"s"`
			}
		)
		var v Inner

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"s":[]}`, string(b))
	})

	t.Run("in struct", func(t *testing.T) {
		type (
			Inner struct {
				S []string `json:"s"`
			}
			Outer struct {
				I Inner `json:"i"`
			}
		)
		v := Outer{}

		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"i":{"s":[]}}`, string(b))
	})

	t.Run("in map", func(t *testing.T) {
		type Inner struct {
			S []string `json:"s"`
		}
		v := map[string]Inner{
			"1": {},
		}
		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"1":{"s":[]}}`, string(b))
	})

	t.Run("in map as interface{}", func(t *testing.T) {
		type Inner struct {
			O string   `json:"o"`
			S []string `json:"s"`
		}
		v := map[string]any{
			"1": any(Inner{O: "foo"}),
		}
		b, err := json.Marshal(niltoempty.Initialize(&v))
		require.NoError(t, err)
		assert.Equal(t, `{"1":{"o":"foo","s":[]}}`, string(b))
	})
}

func TestComplexSlices(t *testing.T) {
	var (
		emptySlice1, emptySlice2 []any
		emptyMap1, emptyMap2     map[string]any
	)
	v := []any{
		nil,
		emptySlice1,
		emptyMap1,
		[2][]any{
			nil,
			{
				[]string{"a"},
				emptySlice2,
				emptyMap2,
			},
		},
	}
	b2, err := json.Marshal(&v)
	require.NoError(t, err)
	assert.Equal(t, `[null,null,null,[null,[["a"],null,null]]]`, string(b2))

	b, err := json.Marshal(niltoempty.Initialize(&v))
	require.NoError(t, err)
	assert.Equal(t, `[null,[],{},[[],[["a"],[],{}]]]`, string(b))
}

func TestCyclic(t *testing.T) {
	t.Run("via pointer", func(t *testing.T) {
		type TC struct {
			P *TC   `json:"p"`
			S []any `json:"s"`
		}

		v := TC{}
		v.P = &v
		_ = niltoempty.Initialize(&v)
	})

	t.Run("via interface", func(t *testing.T) {
		type TC struct {
			S []any `json:"s"`
		}

		var emptySlice []string
		v := TC{make([]any, 2)}
		v.S[0] = v
		v.S[1] = v
		v.S[0] = emptySlice
		_ = niltoempty.Initialize(&v)
	})

	t.Run("via slice and map", func(t *testing.T) {
		type TC struct {
			S  []TC        `json:"s"`
			SI []any       `json:"si"`
			P  *TC         `json:"p"`
			I  any         `json:"i"`
			MP map[int]*TC `json:"mp"`
			MI map[int]any `json:"mi"`
		}

		v := TC{
			S:  make([]TC, 2),
			SI: make([]any, 2),
			MP: make(map[int]*TC),
			MI: make(map[int]any),
		}
		v.SI[1] = &v
		v.P = &v
		v.MP[0] = &v
		v.MI[0] = v
		v.SI[0] = v
		v.S[1] = v
		//_, _ = json.Marshal(&v)
		_ = niltoempty.Initialize(&v)
	})
}

func TestPrivateFields(t *testing.T) {
	t.Run("struct with private fields", func(t *testing.T) {
		type Inner struct {
			Public  []string         `json:"public"`
			private []string         `json:"private"`
			Mixed   map[string][]int `json:"mixed"`
		}

		type Outer struct {
			InnerPtr *Inner `json:"inner_ptr"`
			private  map[string]string
			Public   map[string]interface{} `json:"public"`
		}

		// Create test data with nil slices and maps
		test := Outer{
			InnerPtr: &Inner{
				Public:  nil,
				private: nil,
				Mixed:   nil,
			},
			private: nil,
			Public:  nil,
		}

		// Initialize the struct
		niltoempty.Initialize(&test)

		// Public fields should be initialized
		assert.NotNil(t, test.Public, "Public field should be initialized")
		assert.Empty(t, test.Public, "Public field should be empty map")

		require.NotNil(t, test.InnerPtr, "InnerPtr should remain non-nil")
		assert.NotNil(t, test.InnerPtr.Public, "Nested public field should be initialized")
		assert.Empty(t, test.InnerPtr.Public, "Nested public field should be empty slice")
		assert.NotNil(t, test.InnerPtr.Mixed, "Nested mixed field should be initialized")
		assert.Empty(t, test.InnerPtr.Mixed, "Nested mixed field should be empty map")

		// Private fields should remain nil
		assert.Nil(t, test.private, "Private field should remain nil")
		assert.Nil(t, test.InnerPtr.private, "Nested private field should remain nil")
	})

	t.Run("recursion with mixed private/public fields", func(t *testing.T) {
		type RecursiveStruct struct {
			Public      []map[string]interface{} `json:"public"`
			private     map[string][]interface{}
			RecursiveP  *RecursiveStruct `json:"recursive_p"`
			recursiveP2 *RecursiveStruct
		}

		// Create a recursive structure
		recursive := &RecursiveStruct{
			Public:      nil,
			private:     nil,
			RecursiveP:  &RecursiveStruct{Public: nil, private: nil},
			recursiveP2: &RecursiveStruct{Public: nil, private: nil},
		}

		// This should not panic and should only initialize the public fields
		niltoempty.Initialize(recursive)

		// Check public fields are initialized
		assert.NotNil(t, recursive.Public, "Public field should be initialized")
		assert.Empty(t, recursive.Public, "Public field should be empty slice")
		assert.Nil(t, recursive.private, "Private field should remain nil")

		// Check nested public struct
		require.NotNil(t, recursive.RecursiveP, "RecursiveP should remain non-nil")
		assert.NotNil(t, recursive.RecursiveP.Public, "Nested public field should be initialized")
		assert.Empty(t, recursive.RecursiveP.Public, "Nested public field should be empty slice")
		assert.Nil(t, recursive.RecursiveP.private, "Nested private field should remain nil")

		// Check nested private struct pointer
		require.NotNil(t, recursive.recursiveP2, "recursiveP2 should remain non-nil")
		// While we can reach fields inside unexported pointer fields, we can't modify them
		// because the struct they belong to is not addressable (can't be set)
		assert.Nil(t, recursive.recursiveP2.Public, "Public field in private struct pointer remains nil")
		assert.Nil(t, recursive.recursiveP2.private, "Private field in private struct pointer remains nil")
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("invalid values", func(t *testing.T) {
		// Interface with untyped nil
		m := map[string]interface{}{
			"valid": []string{},
			"nil":   nil,
		}

		// Should not panic
		niltoempty.Initialize(&m)

		// Check that nil values remain nil
		assert.Nil(t, m["nil"], "Untyped nil should remain nil")
		assert.NotNil(t, m["valid"], "Valid value should not be nil")
	})

	t.Run("map with mixed nil values", func(t *testing.T) {
		// A map with various nil values including typed nils
		var nilSlice []int
		var nilMap map[string]int

		m := map[string]interface{}{
			"untypedNil": nil,
			"nilSlice":   nilSlice,
			"nilMap":     nilMap,
			"nilPtr":     (*string)(nil),
		}

		// Should not panic
		niltoempty.Initialize(&m)

		// Check that slices and maps are initialized but pointers and untyped nils are not
		assert.Nil(t, m["untypedNil"], "Untyped nil should remain nil")
		assert.NotNil(t, m["nilSlice"], "nil slice should be initialized")
		assert.Empty(t, m["nilSlice"], "initialized slice should be empty")
		assert.NotNil(t, m["nilMap"], "nil map should be initialized")
		assert.Empty(t, m["nilMap"], "initialized map should be empty")
		assert.Nil(t, m["nilPtr"], "nil pointer should remain nil")
	})

	t.Run("struct with embedded fields", func(t *testing.T) {
		type Embedded struct {
			Slice []string
			Map   map[string]int
		}

		type Container struct {
			Embedded            // Embedded struct
			ExplicitField []int // Regular field
		}

		c := Container{}

		// Should not panic and should initialize all nil slices and maps
		niltoempty.Initialize(&c)

		// Check that all fields are initialized
		assert.NotNil(t, c.Slice, "Embedded Slice should be initialized")
		assert.Empty(t, c.Slice, "Embedded Slice should be empty")
		assert.NotNil(t, c.Map, "Embedded Map should be initialized")
		assert.Empty(t, c.Map, "Embedded Map should be empty")
		assert.NotNil(t, c.ExplicitField, "ExplicitField should be initialized")
		assert.Empty(t, c.ExplicitField, "ExplicitField should be empty")
	})
}

func TestReflectionLimitations(t *testing.T) {
	t.Run("reflection rules with exported vs unexported", func(t *testing.T) {
		type InnerType struct {
			ExportedSlice   []int
			unexportedSlice []int
		}

		type TestStruct struct {
			// Direct exported pointer to a struct with exported and unexported fields
			ExportedPtr *InnerType

			// Unexported pointer to a struct with exported and unexported fields
			unexportedPtr *InnerType

			// Direct exported struct with exported and unexported fields
			ExportedStruct InnerType

			// Unexported struct with exported and unexported fields
			unexportedStruct InnerType
		}

		// Create test data with all nil slices
		test := TestStruct{
			ExportedPtr: &InnerType{
				ExportedSlice:   nil,
				unexportedSlice: nil,
			},
			unexportedPtr: &InnerType{
				ExportedSlice:   nil,
				unexportedSlice: nil,
			},
			ExportedStruct: InnerType{
				ExportedSlice:   nil,
				unexportedSlice: nil,
			},
			unexportedStruct: InnerType{
				ExportedSlice:   nil,
				unexportedSlice: nil,
			},
		}

		// Initialize the struct
		niltoempty.Initialize(&test)

		// Verify behavior for the exported pointer case:
		// - Exported slice inside the struct pointed to by an exported pointer SHOULD be initialized
		assert.NotNil(t, test.ExportedPtr.ExportedSlice,
			"ExportedSlice inside ExportedPtr should be initialized")
		// - Unexported slice inside remains nil even when the parent struct is accessible
		assert.Nil(t, test.ExportedPtr.unexportedSlice,
			"unexportedSlice inside ExportedPtr should remain nil")

		// Verify behavior for the unexported pointer case:
		// - The struct pointed to by an unexported pointer can be accessed
		// - But the exported fields inside cannot be set because the struct is not addressable
		assert.Nil(t, test.unexportedPtr.ExportedSlice,
			"ExportedSlice inside unexportedPtr should remain nil")
		assert.Nil(t, test.unexportedPtr.unexportedSlice,
			"unexportedSlice inside unexportedPtr should remain nil")

		// Verify behavior for the exported struct case:
		// - Exported slice inside an exported struct field SHOULD be initialized
		assert.NotNil(t, test.ExportedStruct.ExportedSlice,
			"ExportedSlice inside ExportedStruct should be initialized")
		// - Unexported slice inside remains nil even when the parent struct is accessible
		assert.Nil(t, test.ExportedStruct.unexportedSlice,
			"unexportedSlice inside ExportedStruct should remain nil")

		// Verify behavior for the unexported struct case:
		// - We can't set any fields (exported or not) inside an unexported struct field
		assert.Nil(t, test.unexportedStruct.ExportedSlice,
			"ExportedSlice inside unexportedStruct should remain nil")
		assert.Nil(t, test.unexportedStruct.unexportedSlice,
			"unexportedSlice inside unexportedStruct should remain nil")
	})
}
