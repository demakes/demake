package models

import (
	"fmt"
	"reflect"
)

type ModelSchema struct {
	Name           string
	Type           reflect.Type
	Fields         []*ModelSchemaField
	RelatedSchemas []*RelatedModelSchema
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
	Struct
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

func SchemaForType(modelType reflect.Type) *ModelSchema {
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

func SchemaFor(model any) *ModelSchema {
	return SchemaForType(reflect.TypeOf(model))
}

func MakeModelSchema(model any) (*ModelSchema, error) {

	modelType := reflect.TypeOf(model)

	// if this is a pointer to a struct, we "unpoint" it first
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model isn't a struct")
	}

	fields := make([]*ModelSchemaField, 0, modelType.NumField())
	related := make([]*RelatedModelSchema, 0)

fieldsLoop:
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tags := ExtractTags(field, "graph")
		pointer := false

		if HasTag(tags, "ignore") {
			// we ignore this field
			continue
		}

		fieldName := ToSnakeCase(field.Name)

		if nameTag, ok := GetTag(tags, "name"); ok {
			fieldName = nameTag.Value
		} else {
			// we look for a JSON name tag
			jsonTags := ExtractTags(field, "json")
			for _, jsonTag := range jsonTags {
				if jsonTag.Flag {
					fieldName = jsonTag.Name
					break
				}
			}
		}

		fieldType := field.Type

		if fieldType.Kind() == reflect.Pointer {
			pointer = true
			fieldType = fieldType.Elem()
		}

		switch fieldType.Kind() {
		case reflect.Map:
			// map
			if fieldType.Key().Kind() != reflect.String {
				break
			}
			mapType := fieldType.Elem()
			mapSchema := SchemaForType(mapType)
			if mapSchema != nil {
				related = append(related, &RelatedModelSchema{
					Type:        Map,
					Name:        fieldName,
					Field:       field.Name,
					Optional:    pointer,
					ModelSchema: mapSchema,
				})
				continue fieldsLoop
			}
		case reflect.Struct:
			// struct
			structSchema := SchemaForType(fieldType)
			if structSchema != nil {
				related = append(related, &RelatedModelSchema{
					Type:        Struct,
					Name:        fieldName,
					Field:       field.Name,
					Optional:    pointer,
					ModelSchema: structSchema,
				})
				continue fieldsLoop
			}
		case reflect.Slice:
			// slice
			sliceType := fieldType.Elem()
			sliceSchema := SchemaForType(sliceType)
			if sliceSchema != nil {
				related = append(related, &RelatedModelSchema{
					Type:        Slice,
					Name:        fieldName,
					Field:       field.Name, // to do: determine field
					Optional:    pointer,
					ModelSchema: sliceSchema,
				})
				continue fieldsLoop
			}
		}

		// this is a regular field
		fields = append(fields, &ModelSchemaField{
			Name:     fieldName,
			Field:    field.Name,
			Type:     fieldType,
			Optional: pointer,
			Tags:     tags,
		})
	}

	return &ModelSchema{
		Type:           modelType,
		Fields:         fields,
		RelatedSchemas: related,
	}, nil
}

func Register[T any](name string) error {
	nt := *new(T)
	if schema, err := MakeModelSchema(nt); err != nil {
		return err
	} else {
		Registry[name] = schema
		return nil
	}
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
