package models

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"reflect"
	"sort"
)

type Hash struct {
	h hash.Hash
}

func MakeHash() *Hash {
	return &Hash{
		h: sha256.New(),
	}
}

func (h *Hash) Add(value any) error {
	return addHash(value, h.h)
}

func (h *Hash) Sum() []byte {
	return h.h.Sum(nil)
}

type CustomHashValue interface {
	HashValue() any
}

var NullValue = fmt.Errorf("null")

func addValue(sourceValue reflect.Value, h hash.Hash) error {

	if sourceValue.IsZero() {
		return NullValue
	}

	sourceType := sourceValue.Type()

	// if the type implements a custom hash value we add this instead of the normal one
	if sourceType.Implements(reflect.TypeOf((*CustomHashValue)(nil)).Elem()) {
		chv := sourceValue.Interface().(CustomHashValue)
		return addHash(chv.HashValue(), h)
	}

	switch sourceType.Kind() {
	case reflect.Slice:
		if sourceValue.Len() == 0 {
			return NullValue
		}
		elemType := sourceType.Elem()
		switch elemType.Kind() {
		case reflect.Uint8: // this is a []byte array
			addHash("bytes", h)
			if _, err := h.Write(sourceValue.Bytes()); err != nil {
				return err
			}
		default: // this is a generic list
			addHash("list", h)
			for i := 0; i < sourceValue.Len(); i++ {
				addHash(fmt.Sprintf("%d", i), h)
				if err := addValue(sourceValue.Index(i), h); err != nil {
					if err == NullValue {
						continue
					}
					return err
				}
			}
		}
	case reflect.Map:
		addHash("map", h)
		if sourceType.Key().Kind() != reflect.String {
			return fmt.Errorf("can only hash string maps")
		}

		stringKeys := make([]string, sourceValue.Len())

		for i, mapKey := range sourceValue.MapKeys() {
			stringKeys[i] = mapKey.String()
		}

		// we sort the string keys
		sort.Strings(stringKeys)

		for _, stringKey := range stringKeys {
			if err := addValue(sourceValue.MapIndex(reflect.ValueOf(stringKey)), h); err != nil {
				if err == NullValue {
					continue
				}
				return err
			} else {
				addHash(stringKey, h)
			}
		}
	case reflect.Struct:

		// we treat structs as equivalent to maps
		addHash("map", h)

		fieldNames := make([]string, sourceType.NumField())
		nameMap := map[string]string{}

		for i := 0; i < sourceType.NumField(); i++ {
			field := sourceType.Field(i)
			tags := ExtractTags(field, "json")

			fieldName := field.Name

			if len(tags) > 0 && tags[0].Flag {
				fieldName = tags[0].Name
			}
			fieldNames[i] = fieldName
			nameMap[fieldName] = field.Name
		}

		sort.Strings(fieldNames)

		for _, fieldName := range fieldNames {
			if err := addValue(sourceValue.FieldByName(nameMap[fieldName]), h); err != nil {
				if err == NullValue {
					continue
				}
				return err
			} else {
				addHash(fieldName, h)
			}
		}
	case reflect.Ptr:
		return addValue(sourceValue.Elem(), h)
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		bs := make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(bs, sourceValue.Int())
		if _, err := h.Write(bs); err != nil {
			return err
		}
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		addHash("int", h)
		bs := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(bs, sourceValue.Uint())
		if _, err := h.Write(bs); err != nil {
			return err
		}
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		addHash("float", h)
		bits := math.Float64bits(sourceValue.Float())
		bs := make([]byte, binary.MaxVarintLen64)
		binary.LittleEndian.PutUint64(bs, bits)
		if _, err := h.Write(bs); err != nil {
			return err
		}
	case reflect.String:
		h.Write([]byte("string"))
		if _, err := h.Write([]byte(sourceValue.String())); err != nil {
			return err
		}
	case reflect.Bool:
		addHash("bool", h)
		if sourceValue.Bool() {
			return addHash(1, h)
		}
		return addHash(0, h)
	case reflect.Interface:
		if sourceValue.IsNil() {
			return NullValue
		} else {
			return addHash(sourceValue.Interface(), h)
		}
	default:
		return fmt.Errorf("unknown type '%v', can't hash", sourceValue.Kind())
	}
	return nil
}

func addHash(source any, h hash.Hash) error {

	sourceValue := reflect.ValueOf(source)

	return addValue(sourceValue, h)

}
