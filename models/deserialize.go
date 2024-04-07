package models

import (
	"fmt"
	"reflect"
)

func Deserialize(node *Node) (any, error) {

	schema, ok := Registry[node.Type]

	if !ok {
		return nil, fmt.Errorf("unknown node type: %s", node.Type)
	}

	modelPtr := reflect.New(schema.Type)
	model := modelPtr.Elem()

	// first, we deserialize the normal data fields
	for _, field := range schema.Fields {
		structField := model.FieldByName(field.Field)
		fieldData, ok := node.JSON.JSON[field.Name]
		if !ok && !field.Optional {
			return nil, fmt.Errorf("field '%s' missing", field.Name)
		}
		dataField := reflect.ValueOf(fieldData)
		dataType := dataField.Type()
		structType := structField.Type()
		if !dataType.AssignableTo(structType) {
			if !dataType.ConvertibleTo(structType) {
				return nil, fmt.Errorf("invalid type (%s expected, but got %s) for field '%s'", structType.Name(), dataType.Name(), field.Name)
			}
			// we convert and set
			structField.Set(dataField.Convert(structType))
		} else {
			structField.Set(dataField)
		}
	}

	for _, relatedSchema := range schema.RelatedSchemas {
		structField := model.FieldByName(relatedSchema.Field)
		structType := structField.Type()
		edges := node.Outgoing.FilterByName(relatedSchema.Name)
		switch relatedSchema.Type {
		case Map:
			mapValue := reflect.MakeMap(structType)
			for _, edge := range edges {
				if model, err := Deserialize(edge.To); err != nil {
					return nil, fmt.Errorf("cannot deserialize related node '%s'(%s): %v", relatedSchema.Name, edge.Key, err)
				} else {
					modelValue := reflect.ValueOf(model)
					if !modelValue.Type().AssignableTo(structType.Elem()) {
						return nil, fmt.Errorf("invalid type: %v vs. %v", modelValue.Type(), structType)
					}
					keyValue := reflect.ValueOf(edge.Key)
					// we set the model value in the map under the given key
					mapValue.SetMapIndex(keyValue, modelValue)
				}
			}
			structField.Set(mapValue)
		case Struct:
			if len(edges) != 1 {
				if len(edges) == 0 && relatedSchema.Optional {
					// this is an optional null value, we skip it
					continue
				}
				return nil, fmt.Errorf("expected exactly one edge, got %d", len(edges))
			}
			if model, err := Deserialize(edges[0].To); err != nil {
				return nil, fmt.Errorf("cannot deserialize related node '%s': %v", relatedSchema.Name, err)
			} else {
				modelValue := reflect.ValueOf(model)
				if structType.Kind() != reflect.Pointer {
					// this isn't a pointer
					modelValue = modelValue.Elem()
				}
				if !modelValue.Type().AssignableTo(structType) {
					return nil, fmt.Errorf("invalid type: %v vs. %v", modelValue.Type(), structType)
				}
				structField.Set(modelValue)
			}
		case Slice:
			// to do: sort the edges!
			for _, edge := range edges {
				if model, err := Deserialize(edge.To); err != nil {
					return nil, fmt.Errorf("cannot deserialize related node '%s'(%d): %v", relatedSchema.Name, edge.Index, err)
				} else {
					modelValue := reflect.ValueOf(model)
					if !modelValue.Type().AssignableTo(structType.Elem()) {
						return nil, fmt.Errorf("invalid type: %v vs. %v", modelValue.Type(), structType)
					}
					// we append the model value to the slice
					structField.Set(reflect.Append(structField, modelValue))
				}
			}
		}
	}

	return modelPtr.Interface(), nil
}
