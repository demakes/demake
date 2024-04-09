package models

import (
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

type Tags []Tag

type Tag struct {
	Name  string
	Value string
	Flag  bool
}

func (t Tags) Has(key string) bool {
	for _, tag := range t {
		if tag.Name == key {
			return true
		}
	}
	return false
}

func (t Tags) Get(key string) (Tag, bool) {
	for _, tag := range t {
		if tag.Name == key {
			return tag, true
		}
	}
	return Tag{}, false
}

func ExtractTags(field reflect.StructField, name string) Tags {
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
