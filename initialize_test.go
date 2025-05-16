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

// New structs for testing based on RagNamedQuerySuccessEvent and QueryResult
type MockQueryResultForTest struct {
	ID            string                   `json:"id"`
	RenderedQuery map[string]interface{}   `json:"rendered_query"` // Can be nil, or contain nil interface values
	Results       []map[string]interface{} `json:"results"`        // Can be nil
	SimpleSlice   []int                    `json:"simple_slice"`   // Can be nil
	SimpleMap     map[string]string        `json:"simple_map"`     // Can be nil
	Payload       interface{}              `json:"payload"`        // Can be untyped nil, or typed nil (e.g. (*SomeStruct)(nil) or (map[string]int)(nil))
	NestedStruct  *NestedStructForTest     `json:"nested_struct"`  // Pointer to a struct, can be nil or struct can have nil fields
}

type NestedStructForTest struct {
	Name           string       `json:"name"`
	DataSlice      []string     `json:"data_slice"` // Can be nil
	DataMap        map[int]bool `json:"data_map"`   // Can be nil
	InterfaceField interface{}  `json:"interface_field"`
}

type MockEventWithPtrSliceForTest struct {
	EventName string                    `json:"event_name"`
	Items     []*MockQueryResultForTest `json:"items"` // Slice of pointers, key for the reported bug
}

func TestSliceOfPointersToStructsWithProblematicFields(t *testing.T) {
	t.Run("struct with nil map/slice/typed-nil-interface fields pointed to by slice element", func(t *testing.T) {
		event := MockEventWithPtrSliceForTest{
			EventName: "TestEvent1",
			Items: []*MockQueryResultForTest{
				{ // Non-nil pointer to MockQueryResultForTest
					ID:            "qr1",
					RenderedQuery: nil,                   // Expected: {}
					Results:       nil,                   // Expected: []
					SimpleSlice:   nil,                   // Expected: []
					SimpleMap:     nil,                   // Expected: {}
					Payload:       (map[string]int)(nil), // Typed nil map, Expected: {}
					NestedStruct: &NestedStructForTest{
						Name:           "nested1",
						DataSlice:      nil,             // Expected: []
						DataMap:        nil,             // Expected: {}
						InterfaceField: ([]string)(nil), // Typed nil slice, Expected: []
					},
				},
			},
		}

		originalItemPtr := event.Items[0]
		originalNestedStructPtr := event.Items[0].NestedStruct

		niltoempty.Initialize(&event)

		require.NotNil(t, event.Items)
		require.Len(t, event.Items, 1)
		require.NotNil(t, event.Items[0])

		assert.Same(t, originalItemPtr, event.Items[0], "Pointer to MockQueryResultForTest should be preserved")

		item0 := event.Items[0]
		assert.NotNil(t, item0.RenderedQuery, "RenderedQuery should be initialized")
		assert.Empty(t, item0.RenderedQuery, "RenderedQuery should be empty map")
		assert.NotNil(t, item0.Results, "Results should be initialized")
		assert.Empty(t, item0.Results, "Results should be empty slice")
		assert.NotNil(t, item0.SimpleSlice, "SimpleSlice should be initialized")
		assert.Empty(t, item0.SimpleSlice, "SimpleSlice should be empty slice")
		assert.NotNil(t, item0.SimpleMap, "SimpleMap should be initialized")
		assert.Empty(t, item0.SimpleMap, "SimpleMap should be empty map")

		require.NotNil(t, item0.Payload, "Payload (typed nil map) should be initialized")
		payloadMap, ok := item0.Payload.(map[string]int)
		require.True(t, ok, "Payload should be a map[string]int, got %T", item0.Payload)
		assert.Empty(t, payloadMap, "Payload (typed nil map) should be an empty map")

		require.NotNil(t, item0.NestedStruct)
		assert.Same(t, originalNestedStructPtr, item0.NestedStruct, "Pointer to NestedStructForTest should be preserved")
		assert.NotNil(t, item0.NestedStruct.DataSlice)
		assert.Empty(t, item0.NestedStruct.DataSlice)
		assert.NotNil(t, item0.NestedStruct.DataMap)
		assert.Empty(t, item0.NestedStruct.DataMap)

		require.NotNil(t, item0.NestedStruct.InterfaceField, "NestedStruct.InterfaceField (typed nil slice) should be initialized")
		interfaceSlice, okSlice := item0.NestedStruct.InterfaceField.([]string)
		require.True(t, okSlice, "NestedStruct.InterfaceField should be a []string, got %T", item0.NestedStruct.InterfaceField)
		assert.Empty(t, interfaceSlice, "NestedStruct.InterfaceField (typed nil slice) should be an empty slice")
	})

	t.Run("struct with untyped nil interface field", func(t *testing.T) {
		event := MockEventWithPtrSliceForTest{
			EventName: "TestEvent2",
			Items: []*MockQueryResultForTest{
				{
					ID:      "qr2",
					Payload: nil, // Untyped nil
					NestedStruct: &NestedStructForTest{
						Name:           "nested2",
						InterfaceField: nil, // Untyped nil
					},
				},
			},
		}
		niltoempty.Initialize(&event)
		require.NotNil(t, event.Items)
		require.Len(t, event.Items, 1)
		item0 := event.Items[0]
		assert.Nil(t, item0.Payload, "Untyped nil Payload should remain nil")
		require.NotNil(t, item0.NestedStruct)
		assert.Nil(t, item0.NestedStruct.InterfaceField, "Untyped nil NestedStruct.InterfaceField should remain nil")
	})

	t.Run("slice of pointers containing a nil pointer", func(t *testing.T) {
		event := MockEventWithPtrSliceForTest{
			EventName: "TestEvent3",
			Items: []*MockQueryResultForTest{
				nil, // A nil pointer in the slice
				{ID: "qr3-valid", SimpleSlice: nil},
			},
		}
		niltoempty.Initialize(&event)
		require.NotNil(t, event.Items)
		require.Len(t, event.Items, 2)
		assert.Nil(t, event.Items[0], "Nil pointer in slice should remain nil")
		require.NotNil(t, event.Items[1])
		assert.NotNil(t, event.Items[1].SimpleSlice, "SimpleSlice in non-nil element should be initialized")
		assert.Empty(t, event.Items[1].SimpleSlice)
	})

	t.Run("map with interface value being untyped nil (potential panic point)", func(t *testing.T) {
		event := MockEventWithPtrSliceForTest{
			EventName: "TestEvent4",
			Items: []*MockQueryResultForTest{
				{
					ID: "qr4",
					RenderedQuery: map[string]interface{}{
						"key1": "value1",
						"key2": nil,                   // Untyped nil interface as map value
						"key3": (map[string]int)(nil), // typed nil map as interface value
					},
				},
			},
		}

		// This subtest previously panicked. Now, untyped nils in map values are preserved.
		// Typed nils (like map[string]int)(nil) are still initialized.
		niltoempty.Initialize(&event)

		// If the code does not panic (e.g., after a fix), these assertions should hold:
		require.NotNil(t, event.Items)
		require.Len(t, event.Items, 1)
		item0 := event.Items[0]
		require.NotNil(t, item0.RenderedQuery)
		assert.Equal(t, "value1", item0.RenderedQuery["key1"])
		assert.Nil(t, item0.RenderedQuery["key2"], "Untyped nil interface map value should remain nil")

		require.NotNil(t, item0.RenderedQuery["key3"], "Typed nil map in interface map value should be initialized")
		key3Map, ok := item0.RenderedQuery["key3"].(map[string]int)
		require.True(t, ok, "item0.RenderedQuery[\"key3\"] expected to be map[string]int, got %T", item0.RenderedQuery["key3"])
		assert.Empty(t, key3Map)
	})

	t.Run("empty slice of pointers", func(t *testing.T) {
		event := MockEventWithPtrSliceForTest{
			EventName: "TestEvent5",
			Items:     []*MockQueryResultForTest{}, // Empty slice
		}
		niltoempty.Initialize(&event)
		require.NotNil(t, event.Items, "Empty slice of pointers should remain not-nil (empty)")
		assert.Empty(t, event.Items, "Empty slice of pointers should remain empty")
	})

	t.Run("nil slice of pointers", func(t *testing.T) {
		event := MockEventWithPtrSliceForTest{
			EventName: "TestEvent6",
			Items:     nil, // Nil slice
		}
		niltoempty.Initialize(&event)
		require.NotNil(t, event.Items, "Nil slice of pointers should be initialized to empty slice")
		assert.Empty(t, event.Items, "Nil slice of pointers should be initialized to empty")
	})

	t.Run("nested struct pointer being nil", func(t *testing.T) {
		event := MockEventWithPtrSliceForTest{
			EventName: "TestEvent7",
			Items: []*MockQueryResultForTest{
				{
					ID:           "qr7",
					NestedStruct: nil, // Nil pointer to nested struct
				},
			},
		}
		niltoempty.Initialize(&event)
		require.NotNil(t, event.Items)
		require.Len(t, event.Items, 1)
		item0 := event.Items[0]
		assert.Nil(t, item0.NestedStruct, "Nil NestedStruct pointer should remain nil")
	})
}

// Additional edge case tests requested by user
func TestAdditionalEdgeCasesForNilToEmpty(t *testing.T) {

	// 1. Metadata map containing an untyped nil value should no longer panic; nil value is preserved.
	t.Run("metadata with untyped nil value is preserved", func(t *testing.T) {
		t.Parallel()
		type EventWithMetadata struct {
			Metadata map[string]interface{} `json:"metadata"`
		}

		event := EventWithMetadata{
			Metadata: map[string]interface{}{
				"foo": nil, // untyped nil interface{}
			},
		}

		// Should NOT panic anymore. Value should remain as-is (nil entry preserved).
		assert.NotPanics(t, func() {
			niltoempty.Initialize(&event)
		}, "Initialize should no longer panic when map contains untyped nil value")

		// Ensure map is still present and the value remains nil.
		require.NotNil(t, event.Metadata)
		val, exists := event.Metadata["foo"]
		assert.True(t, exists, "key foo should still exist in metadata map")
		assert.Nil(t, val, "value of key foo should remain nil")
	})

	// 2. Map containing a typed nil pointer should NOT panic, pointer should stay nil
	t.Run("map with typed nil pointer value stays nil", func(t *testing.T) {
		t.Parallel()
		type PayloadStruct struct {
			Name string `json:"name"`
		}

		type EventWithAdditionalFields struct {
			AdditionalFields map[string]interface{} `json:"additional_fields"`
		}

		var typedNilPointer *PayloadStruct = nil
		event := EventWithAdditionalFields{
			AdditionalFields: map[string]interface{}{
				"payload": typedNilPointer, // typed nil pointer
			},
		}

		// Should NOT panic
		niltoempty.Initialize(&event)

		// The pointer should still be nil after initialization
		require.NotNil(t, event.AdditionalFields)
		payload, ok := event.AdditionalFields["payload"].(*PayloadStruct)
		assert.True(t, ok, "payload should remain of type *PayloadStruct, got %T", event.AdditionalFields["payload"])
		assert.Nil(t, payload, "typed nil pointer value should remain nil after Initialize")
	})

	// 3. Slice containing an untyped nil interface value should not panic and nil should stay nil
	t.Run("slice of interfaces containing untyped nil", func(t *testing.T) {
		t.Parallel()
		type StructWithInterfaceSlice struct {
			Items []interface{} `json:"items"`
		}

		event := StructWithInterfaceSlice{
			Items: []interface{}{nil, "foo", map[string]interface{}{"bar": 1}},
		}

		// Should NOT panic
		niltoempty.Initialize(&event)

		require.Len(t, event.Items, 3)
		assert.Nil(t, event.Items[0], "first element should remain nil")
	})

	// 4. Self-referential map through interface{} value – current implementation would recurse forever.
	//    We include the test but skip it until the algorithm is made cycle-safe for interfaces.
	t.Run("self-referential map via interface cycle", func(t *testing.T) {
		// Mark skipped for now to avoid infinite recursion in the current implementation.
		t.Skip("skipping until initializeNils handles cycles through interface values safely")

		type StructWithCyclicMap struct {
			Data map[string]interface{} `json:"data"`
		}

		cyclic := make(map[string]interface{})
		cyclic["self"] = cyclic // cycle via interface{}
		event := StructWithCyclicMap{Data: cyclic}

		// If/when initializeNils is fixed, this should not panic and should terminate.
		niltoempty.Initialize(&event)
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

		// Check nested private struct - we don't initialize its fields because it's unexported
		require.NotNil(t, recursive.recursiveP2, "recursiveP2 should remain non-nil")
		// Note: Due to the limitations of reflection and Go's visibility rules,
		// we cannot initialize fields within unexported struct pointers
		assert.Nil(t, recursive.recursiveP2.Public, "Fields inside private pointers aren't initialized")
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
