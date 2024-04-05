package models

import (
	"fmt"
	"github.com/gospel-sh/gospel/orm"
	"reflect"
	"regexp"
	"strings"
)

// https://gist.github.com/stoewer/fbe273b711e6a06315d19552dd4d33e6f
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

type Tag struct {
	Name  string
	Value string
	Flag  bool
}

func HasTag(tags []Tag, key string) bool {
	for _, t := range tags {
		if t.Name == key {
			return true
		}
	}
	return false
}

func GetTag(tags []Tag, key string) (Tag, bool) {
	for _, t := range tags {
		if t.Name == key {
			return t, true
		}
	}
	return Tag{}, false
}

func ExtractTags(field reflect.StructField, name string) []Tag {
	tags := make([]Tag, 0)
	if dbValue, ok := field.Tag.Lookup(name); ok {
		strTags := strings.Split(dbValue, ",")
		for _, tag := range strTags {
			kv := strings.Split(dbValue, ":")
			if len(kv) == 1 {
				tags = append(tags, Tag{
					Name:  tag,
					Value: "",
					Flag:  true,
				})
			} else {
				tags = append(tags, Tag{
					Name:  kv[0],
					Value: kv[1],
					Flag:  false,
				})
			}
		}
	}
	return tags
}

func Serialize(model any, db func() orm.DB) (*Node, error) {

	hash := MakeHash()

	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	// if this is a pointer to a struct, we "unpoint" it first
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}

	// we check that this is indeed a struct
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct")
	}

	schema := SchemaFor(model)

	if schema == nil {
		return nil, fmt.Errorf("unknown node type: %T", model)
	}

	node := MakeNode(db)
	data := map[string]any{}

	for _, field := range schema.Fields {
		fieldValue := modelValue.FieldByName(field.Field)
		// we skip zero values
		if fieldValue.IsZero() {
			data[field.Name] = nil
			continue
		}
		// this shouldn't happen, but can...
		if !fieldValue.CanInterface() {
			return nil, fmt.Errorf("cannot interface: %s", field.Field)
		}
		v := fieldValue.Interface()
		if err := hash.Add(v); err != nil {
			return nil, fmt.Errorf("error hashing field %s: %v", field.Field, err)
		}
		data[field.Name] = v
	}

	for _, relatedSchema := range schema.RelatedSchemas {
		fieldValue := modelValue.FieldByName(relatedSchema.Field)

		if fieldValue.IsZero() {
			if !relatedSchema.Optional {
				return nil, fmt.Errorf("related schema '%s' not defined but isn't optional", relatedSchema.Name)
			}
			// we skip this
			continue
		}

		switch relatedSchema.Type {
		case Struct:
			// if this is a pointer to a struct, we "unpoint" it first
			if fieldValue.Kind() == reflect.Pointer {
				fieldValue = fieldValue.Elem()
			}
			if fieldValue.Kind() != reflect.Struct {
				return nil, fmt.Errorf("expected a struct")
			}
			if relatedNode, err := Serialize(fieldValue.Interface(), db); err != nil {
				return nil, fmt.Errorf("cannot serialize related model: %v", err)
			} else {
				edge := MakeEdge(db)
				// we set the type to struct
				edge.Type = int(Struct)
				// we set the edge name
				edge.Name = relatedSchema.Name
				// we link the edge to the nodes
				edge.FromTo(node, relatedNode)

				if err := hash.Add([]any{"edge", edge.Type, "name", edge.Name, "hash", relatedNode.Hash}); err != nil {
					return nil, fmt.Errorf("cannot add edge hash: %v", err)
				}

			}
		case Map:
			if fieldValue.Kind() != reflect.Map {
				return nil, fmt.Errorf("expected a map")
			}
			if fieldValue.Type().Key().Kind() != reflect.String {
				return nil, fmt.Errorf("expected a string key")
			}
			iter := fieldValue.MapRange()
			for iter.Next() {
				mapKey := iter.Key()
				mapValue := iter.Value()
				if mapValue.Kind() == reflect.Pointer {
					mapValue = mapValue.Elem()
				}
				if mapValue.Kind() != reflect.Struct {
					return nil, fmt.Errorf("expected a struct")
				}
				if mapValue.IsZero() || mapKey.IsZero() {
					return nil, fmt.Errorf("found a nil value or key in a map")
				}
				if relatedNode, err := Serialize(mapValue.Interface(), db); err != nil {
					return nil, fmt.Errorf("cannot serialize related model: %v", err)
				} else {
					edge := MakeEdge(db)
					// we set the type to map
					edge.Type = int(Map)
					// we set the edge key to denote it as a map
					edge.Key = mapKey.Interface().(string)
					edge.Name = relatedSchema.Name
					// we link the edge to the nodes
					edge.FromTo(node, relatedNode)

					if err := hash.Add([]any{"edge", edge.Type, "name", edge.Name, "key", edge.Key, "hash", relatedNode.Hash}); err != nil {
						return nil, fmt.Errorf("cannot add edge hash: %v", err)
					}

				}
			}
		case Slice:
			// this is a list of
			if fieldValue.Kind() != reflect.Slice {
				return nil, fmt.Errorf("expected a slice")
			}
			for i := 0; i < fieldValue.Len(); i++ {
				sliceValue := fieldValue.Index(i)
				if sliceValue.Kind() == reflect.Pointer {
					sliceValue = sliceValue.Elem()
				}
				if sliceValue.Kind() != reflect.Struct {
					return nil, fmt.Errorf("expected a struct")
				}
				if sliceValue.IsZero() {
					return nil, fmt.Errorf("found a nil value in a slice")
				}
				if relatedNode, err := Serialize(sliceValue.Interface(), db); err != nil {
					return nil, fmt.Errorf("cannot serialize related model: %v", err)
				} else {
					edge := MakeEdge(db)
					// we set the type to slice
					edge.Type = int(Slice)
					// we set the edge index to denote it as a slice
					edge.Index = i
					edge.Name = relatedSchema.Name
					// we link the edge to the nodes
					edge.FromTo(node, relatedNode)

					if err := hash.Add([]any{"edge", edge.Type, "name", edge.Name, "index", edge.Index, "hash", relatedNode.Hash}); err != nil {
						return nil, fmt.Errorf("cannot add edge hash: %v", err)
					}

				}
			}
		}
	}

	// we generate the node hash
	node.Hash = hash.Sum()
	// we set the node type
	node.Type = schema.Name

	if err := node.JSON.Update(data); err != nil {
		return nil, fmt.Errorf("cannot set data: %v", err)
	}

	return node, nil
}
