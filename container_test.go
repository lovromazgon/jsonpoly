package jsonpoly

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type Animal interface {
	Type() string
	Name() string
}

type Dog struct {
	XName string `json:"name"`
	Breed string `json:"breed"`
}

func (Dog) Type() string {
	return "dog"
}
func (d Dog) Name() string {
	return d.XName
}

type Cat struct {
	XName string `json:"name"`
	Owner string `json:"owner"`
	Color string `json:"color"`
}

func (Cat) Type() string {
	return "cat"
}
func (c Cat) Name() string {
	return c.XName
}

type Pikachu struct{}

func (Pikachu) Type() string { return "pikachu" }
func (Pikachu) Name() string { return "Pikachu" }

type UnknownAnimal struct {
	XType string `json:"-"`
	XName string `json:"name"`
}

func (a UnknownAnimal) Type() string {
	return a.XType
}

func (a UnknownAnimal) Name() string {
	return a.XName
}

var (
	KnownAnimals = map[string]Animal{
		Dog{}.Type():     Dog{},
		Cat{}.Type():     Cat{},
		Pikachu{}.Type(): Pikachu{},
	}
)

type AnimalContainerHelper struct {
	Type string `json:"type"`
}

func (h *AnimalContainerHelper) Get() Animal {
	if a, ok := KnownAnimals[h.Type]; ok {
		return a
	}
	return UnknownAnimal{XType: h.Type}
}

func (h *AnimalContainerHelper) Set(a Animal) {
	h.Type = a.Type()
}

func ExampleContainer_marshal() {
	dog := Dog{
		XName: "Fido",
		Breed: "Golden Retriever",
	}

	var c Container[Animal, *AnimalContainerHelper]
	c.Value = dog

	raw, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(raw))

	// Output:
	// {"type":"dog","name":"Fido","breed":"Golden Retriever"}
}

func ExampleContainer_unmarshal() {
	raw := `{"type":"dog","name":"Fido","breed":"Golden Retriever"}`

	var c Container[Animal, *AnimalContainerHelper]

	err := json.Unmarshal([]byte(raw), &c)
	if err != nil {
		panic(err)
	}

	dog := c.Value.(Dog)

	fmt.Printf("Type: %s\n", dog.Type())
	fmt.Printf("Name: %s\n", dog.XName)
	fmt.Printf("Breed: %s\n", dog.Breed)

	// Output:
	// Type: dog
	// Name: Fido
	// Breed: Golden Retriever
}

func TestContainer_value(t *testing.T) {
	testCases := []struct {
		name string
		have Animal
		want string
	}{
		{
			name: "dog",
			have: Dog{
				XName: "Fido",
				Breed: "Golden Retriever",
			},
			want: `{"type":"dog","name":"Fido","breed":"Golden Retriever"}`,
		},
		{
			name: "cat",
			have: Cat{
				XName: "Whiskers",
				Owner: "Alice",
				Color: "White",
			},
			want: `{"type":"cat","name":"Whiskers","owner":"Alice","color":"White"}`,
		},
		{
			name: "pikachu",
			have: Pikachu{},
			want: `{"type":"pikachu"}`,
		},
		{
			name: "dolphin",
			have: UnknownAnimal{
				XType: "dolphin",
				XName: "Cooper",
			},
			want: `{"type":"dolphin","name":"Cooper"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_marshal", tc.name), func(t *testing.T) {
			c := Container[Animal, *AnimalContainerHelper]{
				Value: tc.have,
			}

			got, err := json.Marshal(c)
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != tc.want {
				t.Fatalf("want %s, got %s", tc.want, string(got))
			}
		})
		t.Run(fmt.Sprintf("%s_unmarshal", tc.name), func(t *testing.T) {
			var c Container[Animal, *AnimalContainerHelper]
			err := json.Unmarshal([]byte(tc.want), &c)
			if err != nil {
				t.Fatal(err)
			}

			got := c.Value
			if got != tc.have {
				t.Fatalf("want %v, got %v", tc.have, got)
			}
		})
	}
}

// AnimalPtrContainerHelper is the same as AnimalContainerHelper, except that it
// returns pointers instead of values in Get.
type AnimalPtrContainerHelper struct {
	Type string `json:"type"`
}

func (h *AnimalPtrContainerHelper) Get() Animal {
	knownAnimals := map[string]Animal{
		Dog{}.Type():     &Dog{},
		Cat{}.Type():     &Cat{},
		Pikachu{}.Type(): &Pikachu{},
	}

	if a, ok := knownAnimals[h.Type]; ok {
		return a
	}
	return &UnknownAnimal{XType: h.Type}
}

func (h *AnimalPtrContainerHelper) Set(a Animal) {
	h.Type = a.Type()
}

func TestContainer_pointer(t *testing.T) {
	testCases := []struct {
		name string
		have Animal
		want string
	}{
		{
			name: "dog",
			have: &Dog{
				XName: "Fido",
				Breed: "Golden Retriever",
			},
			want: `{"type":"dog","name":"Fido","breed":"Golden Retriever"}`,
		},
		{
			name: "cat",
			have: &Cat{
				XName: "Whiskers",
				Owner: "Alice",
				Color: "White",
			},
			want: `{"type":"cat","name":"Whiskers","owner":"Alice","color":"White"}`,
		},
		{
			name: "pikachu",
			have: &Pikachu{},
			want: `{"type":"pikachu"}`,
		},
		{
			name: "dolphin",
			have: &UnknownAnimal{
				XType: "dolphin",
				XName: "Cooper",
			},
			want: `{"type":"dolphin","name":"Cooper"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_marshal", tc.name), func(t *testing.T) {
			c := Container[Animal, *AnimalPtrContainerHelper]{
				Value: tc.have,
			}

			got, err := json.Marshal(c)
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != tc.want {
				t.Fatalf("want %s, got %s", tc.want, string(got))
			}
		})
		t.Run(fmt.Sprintf("%s_unmarshal", tc.name), func(t *testing.T) {
			var c Container[Animal, *AnimalPtrContainerHelper]
			err := json.Unmarshal([]byte(tc.want), &c)
			if err != nil {
				t.Fatal(err)
			}

			// dereference pointers and compare values
			got := reflect.ValueOf(c.Value).Elem().Interface().(Animal)
			have := reflect.ValueOf(tc.have).Elem().Interface().(Animal)

			if got != have {
				t.Fatalf("want %v, got %v", have, got)
			}
		})
	}
}
