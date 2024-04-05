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

	if err := models.Register[Attribute]("attribute"); err != nil {
		t.Fatal(err)
	}

	if err := models.Register[Tag]("tag"); err != nil {
		t.Fatal(err)
	}

	tag := &Tag{
		Type: "p",
		Attributes: []*Attribute{
			&Attribute{
				Name:  "style",
				Value: "font-size:12px",
			},
			&Attribute{
				Name:  "class",
				Value: "bar",
			},
		},
	}

	tagSchema := models.SchemaFor(tag)

	if tagSchema == nil {
		t.Fatalf("expected a schema")
	}

	if len(tagSchema.Fields) != 1 {
		t.Fatalf("expected one regular field")
	}

	node, err := models.Serialize(tag, nil)

	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]any{"type": "p"}

	if !reflect.DeepEqual(node.JSON.JSON, expected) {
		t.Fatalf("data doesn't match: %v vs. %v", node.JSON.JSON, expected)
	}

	if len(node.Outgoing) != 2 {
		t.Fatalf("expected one outgoing edge, got %d", len(node.Outgoing))
	}

	styleEdge := node.Outgoing[0]
	classEdge := node.Outgoing[1]

	if styleEdge.Index != 0 {
		t.Fatalf("expected 0 index")
	}

	if classEdge.Index != 1 {
		t.Fatalf("expected 1 index")
	}

	expected = map[string]any{"name": "style", "value": "font-size:12px"}
	if !reflect.DeepEqual(styleEdge.To.JSON.JSON, expected) {
		t.Fatalf("data doesn't match: %v vs. %v", styleEdge.To.JSON.JSON, expected)
	}

	expected = map[string]any{"name": "class", "value": "bar"}
	if !reflect.DeepEqual(classEdge.To.JSON.JSON, expected) {
		t.Fatalf("data doesn't match: %v vs. %v", styleEdge.To.JSON.JSON, expected)
	}

}
