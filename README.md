# Go Get Set

[![codecov](https://codecov.io/gh/phungvandat/gese/branch/master/graph/badge.svg)](https://codecov.io/gh/phungvandat/gese)
![test status](https://github.com/phungvandat/gese/actions/workflows/test.yml/badge.svg?branch=master)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## Features

- Get: Get the value at `path` from `object`, if the value at `path` not exists or is zero will return the `default value`.

## Usage

- Explain:
  - `from` argument: `string`, `array`, `slice`, `struct` or `map`,
  - `path` argument: `string`, `number`, `array`
  - `defaultVal` argument: anything you want set it to default return value if `path` not exists
  - `isReplaceDefaultVal` argument: `bool`, with `true` value, when `path` exists and `value` at that place is `zero value`, `defaultVal` will replace `zero value`.

```go
func Get(from, path, defaultVal interface{}, isReplaceDefaultVal bool)
```

- Example

```go
package gese

import (
	"github.com/phungvandat/gese"
)

type A struct {
	B B
}

type B struct {
	C *C
}

type C struct {
	D  D
	LD []D
}

type D struct {
	E string
}

func GetExample() {
	// list
	var list = []string{"a", "b", "c"}
	l1 := gese.Get(list, 1)               // b
	l2 := gese.Get(list, 103232, "hello") // hello

	// struct
	var a = A{
		B: B{
			C: &C{
				D: D{
					E: "hello",
				},
				LD: []D{{"e1"}, {"e2"}, {}},
			},
		},
	}

	v0 := gese.Get()  // nil
	v1 := gese.Get(a) // nil

	v2 := gese.Get(a, []string{"B", "C", "D", "E"}) // hello
	v3 := gese.Get(a, "B.C.D.E")                    // hello
	v4 := gese.Get(a, "B.C.D.E.F", 100)             // 100

	v5 := gese.Get(a, []string{"B", "C", "LD", "0", "E"})    // e1
	v6 := gese.Get(a, "B.C.LD.1.E")                          // e2
	v7 := gese.Get(a, []interface{}{"B", "C", "LD", 2, "E"}, "ed", true) // ed

	// map
	var m = map[string]interface{}{
		"A": map[string]interface{}{
			"B": map[string]interface{}{
				"C": map[string]interface{}{
					"D": 100,
				},
			},
		},
	}

	d1 := gese.Get(m, "A.B.C.D")     // 100
	d2 := gese.Get(m, "A.B.C.D.E.F") // nil
	d3 := gese.Get(m, "X.Y.Z", 1000) // 1000
}
```

# Author

**phungvandat**

- [LinkedIn](https://www.linkedin.com/in/phungvandat)
- [Twitter](https://twitter.com/phungvandat97)

## License

Released under the [MIT License](https://github.com/phungvandat/gese/blob/master/LICENSE).
