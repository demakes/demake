package models_test

import (
	"github.com/klaro-org/sites/models"
	"reflect"
	"testing"
)

type Tag struct {
	Type       string       `json:"type"`
	Attributes []*Attribute `json:"attributes"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func TestSerialize(t *testing.T) {

	models.Register[Tag]("tag")
	models.Register[Attribute]("attribute")

	tag := &Tag{
		Type: "p",
		Attributes: []*Attribute{
			&Attribute{
				Name:  "style",
				Value: "font-size:12px",
			},
		},
	}

	node, err := models.Serialize(tag)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(node.JSON.Get(), map[string]any{"type": "p"}) {
		t.Fatalf("data doesn't match")
	}
}
