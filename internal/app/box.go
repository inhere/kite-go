package app

import (
	"github.com/gookit/goutil"
)

// global object container box
var box = map[string]any{}

// Set object by name
func Set[T any](name string, obj T) {
	box[name] = obj
}

// Add object by name
func Add[T any](name string, obj T) {
	box[name] = obj
}

// Has object by name
func Has(name string) bool {
	_, ok := box[name]
	return ok
}

// GetAny value from box
func GetAny(name string) any {
	return box[name]
}

// Get object by name, if not exists will panic
func Get[T any](name string) T {
	obj, ok := box[name]
	if !ok {
		goutil.Panicf("object %q not exists in box", name)
	}
	return obj.(T)
}

// Lookup object by name
func Lookup[T any](name string) (v T, ok bool) {
	if obj, ok := box[name]; ok {
		return obj.(T), ok
	}
	return
}

// Names object name list from box
func Names() []string {
	names := make([]string, 0, len(box))
	for name := range box {
		names = append(names, name)
	}
	return names
}
