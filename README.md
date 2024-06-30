# JSONPoly

[![License](https://img.shields.io/github/license/lovromazgon/jsonpoly)](https://github.com/ConduitIO/conduit/blob/main/LICENSE)
[![Test](https://github.com/lovromazgon/jsonpoly/actions/workflows/test.yml/badge.svg)](https://github.com/lovromazgon/jsonpoly/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovromazgon/jsonpoly)](https://goreportcard.com/report/github.com/lovromazgon/jsonpoly)
[![Go Reference](https://pkg.go.dev/badge/github.com/lovromazgon/jsonpoly.svg)](https://pkg.go.dev/github.com/lovromazgon/jsonpoly)

Utilities for marshalling and unmarshalling polymorphic JSON objects in Go
without generating code.

## Usage

```
go get github.com/lovromazgon/jsonpoly@latest
```

Say that you have an interface `Shape` and two structs `Triangle` and `Square`
that implement it. The structs have a method `Kind` that returns the name of
the shape.

```go
package shapes

type Shape interface {
	Kind() string
}

func (Triangle) Kind() string { return "triangle" }
func (Square) Kind() string   { return "square" }

type Square struct {
	TopLeft [2]int `json:"top-left"`
	Width   int    `json:"width"`
}

type Triangle struct {
	P0 [2]int `json:"p0"`
	P1 [2]int `json:"p1"`
	P2 [2]int `json:"p2"`
}
```

You need to define a type that implements the `jsonpoly.Helper` interface and
can marshal and unmarshal the field(s) used to determine the type of the object.
In this case, the field is `kind`. You also need to define a map that maps the
values of the field to the types.

```go
var knownShapes = map[string]Shape{
	Triangle{}.Kind(): Triangle{},
	Square{}.Kind():   Square{},
}

type ShapeJSONHelper struct {
	Kind string `json:"kind"`
}

func (h *ShapeJSONHelper) Get() Shape {
	return knownShapes[h.Kind]
}

func (h *ShapeJSONHelper) Set(s Shape) {
	h.Kind = s.Kind()
}
```

Now you can marshal and unmarshal polymorphic JSON objects using `jsonpoly.Container`.

```go
inputShape := Square{TopLeft: [2]int{1, 2}, Width: 4}

var c jsonpoly.Container[Shape, *ShapeJSONHelper]
c.Value = inputShape

b, err := json.Marshal(c)
fmt.Println(string(b)) // {"kind":"square","top-left":[1,2],"width":4}

c.Value = nil // reset before unmarshalling
err = json.Unmarshal(b, &c)
fmt.Printf("%T\n", c.Value) // shapes.Square
```

Also check out the
[marshalling](https://pkg.go.dev/github.com/lovromazgon/jsonpoly#example-Container-Marshal)
and [unmarshalling](https://pkg.go.dev/github.com/lovromazgon/jsonpoly#example-Container-Unmarshal)
examples on the package documentation.

## FAQ

### How is this different than [`github.com/polyfloyd/gopolyjson`](https://github.com/polyfloyd/gopolyjson)?

`gopolyjson` is a great package, but it has its limitations. Here's a list of
differences that can help you determine what package to use:

- `gopolyjson` requires you to add a private method to your interface without
  parameters or return arguments. As a consequence, you have to put all types
  that implement the interface in the same package. `jsonpoly` does not require
  you to add any methods to your types.
- `gopolyjson` requires you to generate code for each type you want to serialize.
  Since the generated code adds methods to the types, you can not generate the
  code for types from external packages. `jsonpoly` works without generating code.
- Because `gopolyjson` uses generated code, it can be faster than `jsonpoly`.
- `gopolyjson` only supports a single field at the root of the JSON to determine
  the type of the object, while `jsonpoly` supports multiple fields.
- `gopolyjson` does not handle unknown types which can be an issue with
  backwards compatibility. `jsonpoly` can handle unknown types by having a
  "catch-all" type.

### How can I handle unknown types?

If you want to handle unknown types, you can define a "catch-all" type. The type
should be returned by the `Get` method of the `jsonpoly.Helper` implementation
whenever the type of the object is not recognized.

Keep in mind that the field used to determine the type of the object should be
marked with the `json:"-"` tag, as it is normally handled by the helper. Not
doing so will result in duplicating the field.

```go
type Unknown struct {
    XKind string `json:"-"`
    json.RawMessage // Store the raw json if needed.
}

func (u Unknown) Kind() string { return u.XKind }

type ShapeJSONHelper struct {
    Kind string `json:"kind"`
}

func (h *ShapeJSONHelper) Get() Shape {
    s, ok := knownShapes[h.Kind]
    if !ok {
        return Unknown{XKind: h.Kind}
    }
    return s
}
```

### Can I use multiple fields to determine the type of the object?

Yes, you can use any number of fields to determine the type of the object. You
just need to define a struct that contains all the fields and implements the
`jsonpoly.Helper` interface.

For more information on how to do this, check the [`example`](./example) directory.