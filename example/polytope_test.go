package example

import (
	"encoding/json"
	"fmt"

	"github.com/lovromazgon/jsonpoly"
)

func ExamplePolytope() {
	inputPolytope := Square{TopLeft: [2]int{1, 2}, Width: 4}

	c := jsonpoly.Container[Polytope, *PolytopeJSONHelper]{Value: inputPolytope}

	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b) // {"kind":"hypercube","dimension":2,"top-left":[1,2],"width":4}

	c.Value = nil
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T\n", c.Value) // examples.Square

	// Output:
	// {"kind":"hypercube","dimension":2,"top-left":[1,2],"width":4}
	// examples.Square
}
