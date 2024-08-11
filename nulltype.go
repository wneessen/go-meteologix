package meteologix

import "encoding/json"

// Variable is a generic variable type that can be null.
type Variable[T any] struct {
	value  T
	notNil bool
}

// Get the value of the Variable
func (f *Variable[T]) Get() T {
	return f.value
}

// NotNil returns true when a Variable is not nil
func (f *Variable[T]) NotNil() bool {
	return f.notNil
}

// IsNil returns true when a Variable is nil
func (f *Variable[T]) IsNil() bool {
	return !f.notNil
}

// Reset resets the value to the Variable to a zero value and sets it to be nil
func (f *Variable[T]) Reset() {
	var newVal T
	f.value = newVal
	f.notNil = false
}

// NilInt is an int type that can be nil
type NilInt = Variable[int]

// NilFloat64 is an float64 type that can be nil
type NilFloat64 = Variable[float64]

// NilString is a string type that can be nil
type NilString = Variable[string]

// UnmarshalJSON interprets the generic Nil types and sets the value and notnil of the type
func (f *Variable[T]) UnmarshalJSON(data []byte) error {
	if string(data) != "null" {
		f.value = *new(T)
		f.notNil = true
		return json.Unmarshal(data, &f.value)
	}
	return nil
}
