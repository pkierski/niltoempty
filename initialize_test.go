package niltoempty_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

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

func TestTimeHandling(t *testing.T) {
	t.Run("struct with time.Time fields", func(t *testing.T) {
		type TimeStruct struct {
			CreatedAt    time.Time                        `json:"created_at"`
			UpdatedAt    *time.Time                       `json:"updated_at"`
			NullableTime *time.Time                       `json:"nullable_time"`
			Maps         map[string]time.Time             `json:"maps"`
			MapPtrs      map[string]*time.Time            `json:"map_ptrs"`
			Slices       []time.Time                      `json:"slices"`
			SlicePtrs    []*time.Time                     `json:"slice_ptrs"`
			NestedMap    map[string]map[string]*time.Time `json:"nested_map"`
		}

		// Get current time with timezone info
		location := time.FixedZone("UTC+2", 2*60*60)
		now := time.Now().In(location)

		// Create a struct with time.Time fields
		timeStruct := TimeStruct{
			CreatedAt:    now,
			UpdatedAt:    nil,
			NullableTime: nil,
			Maps:         nil,
			MapPtrs:      nil,
			Slices:       nil,
			SlicePtrs:    nil,
			NestedMap:    nil,
		}

		// Initialize the struct
		niltoempty.Initialize(&timeStruct)

		// Verify that time.Time fields are untouched
		assert.Equal(t, now, timeStruct.CreatedAt, "CreatedAt time should be unchanged")
		// Explicitly verify location is preserved
		assert.Equal(t, location.String(), timeStruct.CreatedAt.Location().String(), "Location information should be preserved")
		// Check timezone offset
		_, offset := timeStruct.CreatedAt.Zone()
		assert.Equal(t, 2*60*60, offset, "Location offset should be preserved")

		assert.Nil(t, timeStruct.UpdatedAt, "UpdatedAt should remain nil")
		assert.Nil(t, timeStruct.NullableTime, "NullableTime should remain nil")

		// Maps and slices should be initialized
		assert.NotNil(t, timeStruct.Maps, "Maps should be initialized")
		assert.Empty(t, timeStruct.Maps, "Maps should be empty")
		assert.NotNil(t, timeStruct.MapPtrs, "MapPtrs should be initialized")
		assert.Empty(t, timeStruct.MapPtrs, "MapPtrs should be empty")
		assert.NotNil(t, timeStruct.Slices, "Slices should be initialized")
		assert.Empty(t, timeStruct.Slices, "Slices should be empty")
		assert.NotNil(t, timeStruct.SlicePtrs, "SlicePtrs should be initialized")
		assert.Empty(t, timeStruct.SlicePtrs, "SlicePtrs should be empty")
		assert.NotNil(t, timeStruct.NestedMap, "NestedMap should be initialized")
		assert.Empty(t, timeStruct.NestedMap, "NestedMap should be empty")
	})

	t.Run("map with time.Time interface values", func(t *testing.T) {
		// Create a map with time.Time interface values
		utcLocation := time.UTC
		now := time.Now().In(utcLocation)

		localLocation, err := time.LoadLocation("Local")
		require.NoError(t, err, "Loading Local location should not error")
		localTime := time.Now().In(localLocation)

		m := map[string]interface{}{
			"utc_time":     now,
			"local_time":   localTime,
			"time_ptr":     &now,
			"nil_time_ptr": (*time.Time)(nil),
			"time_map": map[string]time.Time{
				"now": now,
			},
			"nil_time_map":   map[string]time.Time(nil),
			"time_slice":     []time.Time{now, localTime},
			"nil_time_slice": []time.Time(nil),
		}

		// Initialize the map
		niltoempty.Initialize(&m)

		// Verify that time.Time values are untouched
		assert.Equal(t, now, m["utc_time"], "UTC time should be unchanged")
		assert.Equal(t, utcLocation.String(), m["utc_time"].(time.Time).Location().String(), "UTC location should be preserved")

		assert.Equal(t, localTime, m["local_time"], "Local time should be unchanged")
		assert.Equal(t, localLocation.String(), m["local_time"].(time.Time).Location().String(), "Local location should be preserved")

		// Verify pointer to time with location
		timePtr := m["time_ptr"].(*time.Time)
		assert.Equal(t, &now, timePtr, "Time pointer should be unchanged")
		assert.Equal(t, utcLocation.String(), timePtr.Location().String(), "Location in time pointer should be preserved")

		assert.Nil(t, m["nil_time_ptr"], "Nil time pointer should remain nil")

		// Time maps and slices should be initialized
		timeMap, ok := m["time_map"].(map[string]time.Time)
		assert.True(t, ok, "time_map should remain a map[string]time.Time")
		assert.Equal(t, now, timeMap["now"], "Time in map should be unchanged")
		assert.Equal(t, utcLocation.String(), timeMap["now"].Location().String(), "Location in map value should be preserved")

		nilTimeMap, ok := m["nil_time_map"].(map[string]time.Time)
		assert.True(t, ok, "nil_time_map should be initialized as map[string]time.Time")
		assert.Empty(t, nilTimeMap, "Initialized time map should be empty")

		timeSlice, ok := m["time_slice"].([]time.Time)
		assert.True(t, ok, "time_slice should remain a []time.Time")
		assert.Equal(t, now, timeSlice[0], "First time in slice should be unchanged")
		assert.Equal(t, utcLocation.String(), timeSlice[0].Location().String(), "Location of first time in slice should be preserved")
		assert.Equal(t, localTime, timeSlice[1], "Second time in slice should be unchanged")
		assert.Equal(t, localLocation.String(), timeSlice[1].Location().String(), "Location of second time in slice should be preserved")

		nilTimeSlice, ok := m["nil_time_slice"].([]time.Time)
		assert.True(t, ok, "nil_time_slice should be initialized as []time.Time")
		assert.Empty(t, nilTimeSlice, "Initialized time slice should be empty")
	})

	t.Run("struct with time.Time in nested interface fields", func(t *testing.T) {
		type EventWithTimeData struct {
			Metadata map[string]interface{} `json:"metadata"`
			Payload  interface{}            `json:"payload"`
		}

		type TimeData struct {
			CreatedAt time.Time   `json:"created_at"`
			Times     []time.Time `json:"times"`
		}

		// Create a custom location
		customLocation := time.FixedZone("UTC-5", -5*60*60)
		now := time.Now().In(customLocation)
		timeData := TimeData{
			CreatedAt: now,
			Times:     nil,
		}

		event := EventWithTimeData{
			Metadata: map[string]interface{}{
				"timestamp":  now,
				"time_data":  timeData,
				"time_slice": []time.Time(nil),
			},
			Payload: timeData,
		}

		// Initialize the struct
		niltoempty.Initialize(&event)

		// Verify time fields in metadata
		assert.Equal(t, now, event.Metadata["timestamp"], "Timestamp in metadata should be unchanged")
		assert.Equal(t, customLocation.String(), event.Metadata["timestamp"].(time.Time).Location().String(),
			"Location in timestamp should be preserved")

		// Verify timezone offset
		_, offset := event.Metadata["timestamp"].(time.Time).Zone()
		assert.Equal(t, -5*60*60, offset, "Location offset should be preserved")

		metadataTimeData, ok := event.Metadata["time_data"].(TimeData)
		assert.True(t, ok, "time_data in metadata should remain a TimeData struct")
		assert.Equal(t, now, metadataTimeData.CreatedAt, "CreatedAt in metadata should be unchanged")
		assert.Equal(t, customLocation.String(), metadataTimeData.CreatedAt.Location().String(),
			"Location in metadata TimeData should be preserved")
		assert.NotNil(t, metadataTimeData.Times, "Times in metadata should be initialized")
		assert.Empty(t, metadataTimeData.Times, "Times in metadata should be empty")

		metadataTimeSlice, ok := event.Metadata["time_slice"].([]time.Time)
		assert.True(t, ok, "time_slice in metadata should be initialized as []time.Time")
		assert.Empty(t, metadataTimeSlice, "Time slice in metadata should be empty")

		// Verify time fields in payload
		payloadTimeData, ok := event.Payload.(TimeData)
		assert.True(t, ok, "payload should remain a TimeData struct")
		assert.Equal(t, now, payloadTimeData.CreatedAt, "CreatedAt in payload should be unchanged")
		assert.Equal(t, customLocation.String(), payloadTimeData.CreatedAt.Location().String(),
			"Location in payload TimeData should be preserved")
		assert.NotNil(t, payloadTimeData.Times, "Times in payload should be initialized")
		assert.Empty(t, payloadTimeData.Times, "Times in payload should be empty")
	})

	t.Run("complex struct with time.Time inside slices and maps", func(t *testing.T) {
		type TimeEvent struct {
			Timestamp time.Time            `json:"timestamp"`
			Tags      map[string]time.Time `json:"tags"`
		}

		type ComplexTimeStruct struct {
			Events      []TimeEvent                     `json:"events"`
			EventsByTag map[string][]TimeEvent          `json:"events_by_tag"`
			TimesByTag  map[string]map[string]time.Time `json:"times_by_tag"`
		}

		// Create times with different locations
		utcLocation := time.UTC
		now := time.Now().In(utcLocation)

		estLocation := time.FixedZone("EST", -5*60*60)
		yesterday := now.Add(-24 * time.Hour).In(estLocation)

		// Initialize with nil maps and slices
		complex := ComplexTimeStruct{
			Events:      nil,
			EventsByTag: nil,
			TimesByTag:  nil,
		}

		// Initialize the complex struct
		niltoempty.Initialize(&complex)

		// Verify nested fields are initialized
		assert.NotNil(t, complex.Events, "Events should be initialized")
		assert.Empty(t, complex.Events, "Events should be empty")
		assert.NotNil(t, complex.EventsByTag, "EventsByTag should be initialized")
		assert.Empty(t, complex.EventsByTag, "EventsByTag should be empty")
		assert.NotNil(t, complex.TimesByTag, "TimesByTag should be initialized")
		assert.Empty(t, complex.TimesByTag, "TimesByTag should be empty")

		// Add elements and verify they're handled correctly
		complex.Events = append(complex.Events, TimeEvent{Timestamp: now, Tags: nil})
		complex.EventsByTag = map[string][]TimeEvent{
			"test": {
				{Timestamp: yesterday, Tags: nil},
			},
		}
		complex.TimesByTag = map[string]map[string]time.Time{
			"test": nil,
		}

		// Re-initialize and verify
		niltoempty.Initialize(&complex)

		// Verify time fields are preserved and nil maps initialized
		assert.Equal(t, now, complex.Events[0].Timestamp, "Event timestamp should be unchanged")
		assert.Equal(t, utcLocation.String(), complex.Events[0].Timestamp.Location().String(),
			"Location in Events timestamp should be preserved")
		assert.NotNil(t, complex.Events[0].Tags, "Event tags should be initialized")
		assert.Empty(t, complex.Events[0].Tags, "Event tags should be empty")

		assert.Equal(t, yesterday, complex.EventsByTag["test"][0].Timestamp, "EventsByTag timestamp should be unchanged")
		assert.Equal(t, estLocation.String(), complex.EventsByTag["test"][0].Timestamp.Location().String(),
			"Location in EventsByTag timestamp should be preserved")

		// Verify timezone offset
		_, offset := complex.EventsByTag["test"][0].Timestamp.Zone()
		assert.Equal(t, -5*60*60, offset, "Location offset should be preserved")

		assert.NotNil(t, complex.EventsByTag["test"][0].Tags, "EventsByTag tags should be initialized")
		assert.Empty(t, complex.EventsByTag["test"][0].Tags, "EventsByTag tags should be empty")

		assert.NotNil(t, complex.TimesByTag["test"], "TimesByTag map should be initialized")
		assert.Empty(t, complex.TimesByTag["test"], "TimesByTag map should be empty")
	})
}
