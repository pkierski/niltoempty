# niltoempty
[![Go Reference](https://pkg.go.dev/badge/github.com/pkierski/niltoempty.svg)](https://pkg.go.dev/github.com/pkierski/niltoempty)
[![Go Report Card](https://goreportcard.com/badge/github.com/pkierski/niltoempty)](https://goreportcard.com/report/github.com/pkierski/niltoempty)

Recursively initializes all nil maps and slices in a given object, so [json.Marshal()](https://pkg.go.dev/encoding/json#Marshal) serializes they as empty object {} or array [] instead of null.

This is more complete solution based on the idea from [nilslice](https://github.com/golang-cz/nilslice). It works not only for nil slices but also for nil maps. 

```go
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
```
