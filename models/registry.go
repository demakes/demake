package models

import (
	"fmt"
	"reflect"
)

type ModelSchema struct {
	Name                string
	Type                reflect.Type
	Fields              []*ModelSchemaField
	RelatedModelSchemas []*RelatedModelSchema
}

// Serializes a given model into a node structure
func (m *ModelSchema) Serialize(model any) *Node {
	/*
		- Create a node
		- Serialize all fields into the nodes JSON data
		- Serialize all related models into nodes
		- Create outgoing edges for the related models
		- Combine the outgoing edges and node
	*/
	return nil
}

type Relation int

const (
	Map Relation = iota
	Slice
	Literal
)

// describes a related model of a given model
// contains information about the foreign key
type RelatedModelSchema struct {
	Type        Relation
	Name        string
	Field       string
	Optional    bool
	ModelSchema *ModelSchema
}

type ModelSchemaField struct {
	Name     string
	Field    string
	Optional bool
	Type     reflect.Type
	Tags     []Tag
}

// maps a given type name to a struct
var Registry = map[string]*ModelSchema{}

func SchemaFor(model any) *ModelSchema {

	modelType := reflect.TypeOf(model)

	// if this is a pointer to a struct, we "unpoint" it first
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}

	for _, t := range Registry {
		if modelType.AssignableTo(t.Type) {
			return t
		}
	}

	return nil
}

func MakeModelSchema(model any) *ModelSchema {

	modelType := reflect.TypeOf(model)

	// if this is a pointer to a struct, we "unpoint" it first
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}

	return &ModelSchema{
		Type: modelType,
	}
}

func Register[T any](name string) {
	nt := *new(T)
	Registry[name] = MakeModelSchema(nt)
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
