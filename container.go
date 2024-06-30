package jsonpoly

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotJSONObject = errors.New("not a JSON object")
)

// Container is a generic struct that can be used to unmarshal polymorphic JSON
// objects into a specific type based on a key. It is using the Helper interface
// to determine the type of the object and to create a new instance of the
// unmarshalled object.
type Container[V any, H Helper[V]] struct {
	Value V
}

// Helper is an interface that must be implemented by the user to
// provide the necessary methods to create and set the value of the object based
// on the key. The struct implementing this interface should be a pointer type
// and should contain public fields annotated with JSON tags that match the keys
// in the JSON object.
type Helper[V any] interface {
	Get() V
	Set(V)
}

func (c *Container[V, H]) UnmarshalJSON(b []byte) error {
	var helper H
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}

	v := helper.Get()

	// Check if the value is a pointer of a value. If it's a pointer, we use it
	// as is. If it's a value, we create a pointer to it for the unmarshalling
	// to work and store the underlying value in the 'Value' field.
	val := reflect.ValueOf(v)
	if !val.IsValid() {
		// Apparently this is an unknown type, marshal the helper to represent
		// the type and include it in the error message. We can safely ignore
		// the error, since the type was already unmarshalled successfully.
		b, _ := json.Marshal(helper)
		return fmt.Errorf("unknown type %v", string(b))
	}

	var ptrVal reflect.Value
	if val.Kind() != reflect.Ptr {
		// Create a new pointer type based on the type of 'v'.
		ptrType := reflect.PointerTo(val.Type())
		// Allocate a new object of this pointer type.
		ptrVal = reflect.New(ptrType.Elem())
		// Set the newly allocated object to the value of 'v'.
		ptrVal.Elem().Set(val)
		// Now 'ptrVal' is a reflect.Value of type '*V' which can be used as a pointer.
		v = ptrVal.Interface().(V)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	if ptrVal.IsValid() {
		// If we used a pointer, we need to get the underlying value.
		c.Value = ptrVal.Elem().Interface().(V)
	} else {
		// If we used the value directly, we store it in the 'Value' field.
		c.Value = v
	}

	return nil
}

func (c Container[V, H]) MarshalJSON() ([]byte, error) {
	helper := reflect.New(reflect.TypeFor[H]().Elem()).Interface().(H)
	helper.Set(c.Value)

	jsonHelper, err := json.Marshal(helper)
	if err != nil {
		return nil, err
	}

	jsonValue, err := json.Marshal(c.Value)
	if err != nil {
		return nil, err
	}

	return mergeJSONObjects(jsonHelper, jsonValue)
}

func mergeJSONObjects(o1, o2 []byte) ([]byte, error) {
	if !isJSONObject(o1) || !isJSONObject(o2) {
		return nil, ErrNotJSONObject
	}

	// We know this is only used internally, we can manipulate the slices.
	// We append the second object to the first one, replacing the closing
	// object bracket with a comma.
	o2[0] = ','
	return append(o1[:len(o1)-1], o2...), nil
}

func isJSONObject(o []byte) bool {
	if len(o) == 0 {
		return false
	}
	return o[0] == '{' && o[len(o)-1] == '}'
}
