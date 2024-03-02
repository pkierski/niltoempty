package niltoempty_test

import (
	"encoding/json"
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
		[2][]any{nil, {[]string{"a"}, emptySlice2, emptyMap2}},
	}
	b2, err := json.Marshal(&v)
	require.NoError(t, err)
	assert.Equal(t, `[null,null,null,[null,[["a"],null,null]]]`, string(b2))

	b, err := json.Marshal(niltoempty.Initialize(&v))
	require.NoError(t, err)
	assert.Equal(t, `[null,[],{},[[],[["a"],[],{}]]]`, string(b))
}
