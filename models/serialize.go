package models

import (
	"fmt"
	"reflect"
	"strings"
)

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

func ExtractTags(field reflect.StructField) []Tag {
	tags := make([]Tag, 0)
	if dbValue, ok := field.Tag.Lookup("graph"); ok {
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

func Serialize(model any) (*Node, error) {
	/*
		- We go through all of the fields of the model
		- If it's a list or a map, we check if the value type is a mapped model
		- If so, we serialize it and create edges for the mapped models
		- If not, we add it to the data of the node
		- We hash the node
	*/

	modelType := reflect.TypeOf(model)
	modelValue := reflect.ValueOf(model)

	// if this is a pointer to a struct, we "unpoint" it first
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
		modelValue = modelValue.Elem()
	}

	schema := SchemaFor(model)

	if schema == nil {
		return nil, fmt.Errorf("unknown node type: %T", model)
	}

	// we check that this is indeed a struct
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct")
	}

	node := MakeNode(nil)
	return node, nil
}
