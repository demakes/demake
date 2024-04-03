package models

import (
	"fmt"
)

// maps a given type name to a struct
var Registry map[string]any = map[string]any{}

func Register[T any](name string) {
	nt := new(T)
	Registry[name] = nt
}

func GetType(name string) any {
	return Registry[name]
}

type MyNode struct {
	Foo string
}

func init() {
	// we register the MyNode type for the 'myNode' type
	Register[MyNode]("myNode")
	// GetType can be used to reflectively create a new struct for the type
	fmt.Printf("Type: %T\n", GetType("myNode"))
}
