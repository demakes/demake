package models_test

import (
	"encoding/hex"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites"
	"github.com/klaro-org/sites/models"
	kt "github.com/klaro-org/sites/testing"
	"reflect"
	"testing"
)

type Tag struct {
	Type       string       `json:"type"`
	Attributes []*Attribute `json:"attributes"`
}

type Attribute struct {
	Name   string            `json:"name"`
	Value  string            `json:"value"`
	Labels map[string]*Label `json:"labels"`
}

type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func TestSerialize(t *testing.T) {

	settings, err := sites.LoadSettings()

	if err != nil {
		t.Fatal(err)
	}

	db, err := kt.DB(settings)

	if err != nil {
		t.Fatal(err)
	}

	// we first register the label model
	if err := models.Register[Label]("label"); err != nil {
		t.Fatal(err)
	}

	// then we register the attribute model
	if err := models.Register[Attribute]("attribute"); err != nil {
		t.Fatal(err)
	}

	// then we register the tag model
	if err := models.Register[Tag]("tag"); err != nil {
		t.Fatal(err)
	}

	tag := &Tag{
		Type: "p",
		Attributes: []*Attribute{
			&Attribute{
				Name:   "style",
				Value:  "font-size:12px",
				Labels: map[string]*Label{},
			},
			&Attribute{
				Name:  "class",
				Value: "bar",
				Labels: map[string]*Label{
					"test": {
						Name:  "foo",
						Value: "bar",
					},
				},
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

	node, err := models.Serialize(tag, func() orm.DB { return db })

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

	h := "bbe85f36e4f6d285e0f0cb34a893a63b2478a5abb45db160ad52e6f60801e0f2"
	if hex.EncodeToString(classEdge.To.Hash) != h {
		t.Fatalf("invalid hash, expected '%s', got '%s'", h, hex.EncodeToString(classEdge.To.Hash))
	}

	// we store the node a first time
	if err := node.SaveTree(); err != nil {
		t.Fatalf("cannot store node")
	}

	// we modify the second attribute
	tag.Attributes[1].Name = "classes"

	newNode, err := models.Serialize(tag, func() orm.DB { return db })

	if err != nil {
		t.Fatal(err)
	}

	// we store the node a second time
	if err := newNode.SaveTree(); err != nil {
		t.Fatalf("cannot store new node a second time")
	}

	if newNode.Outgoing[0].To.ID != node.Outgoing[0].To.ID {
		t.Fatalf("expected IDs of first attribute to match")
	}

	if newNode.Outgoing[1].To.ID == node.Outgoing[1].To.ID {
		t.Fatalf("expected IDs of second attribute to diverge")
	}

	// we check the number of nodes in the database
	nodes, err := orm.Objects[models.Node](func() orm.DB { return db }, map[string]any{})

	if err != nil {
		t.Fatal(err)
	}

	// we expect 6 nodes as the initial tree has 4 and the new tree replaces
	// 2 nodes with modified ones...
	if len(nodes) != 6 {
		t.Fatalf("expected 6 nodes, got %d", len(nodes))
	}

	// we check the number of nodes in the database
	edges, err := orm.Objects[models.Edge](func() orm.DB { return db }, map[string]any{})

	if err != nil {
		t.Fatal(err)
	}

	// we expect 6 edges as each graph has 3
	if len(edges) != 6 {
		t.Fatalf("expected 6 edges, got %d", len(edges))
	}

}
