package models_test

import (
	"encoding/hex"
	"encoding/json"
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

	dbf := func() orm.DB { return db }

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
				Name:  "style",
				Value: "font-size:12px",
				Labels: map[string]*Label{
					"test": {
						Name:  "fooz",
						Value: "bar",
					},
					"baz": {
						Name:  "foo",
						Value: "bar",
					},
				},
			},
			&Attribute{
				Name:  "class",
				Value: "bar",
				Labels: map[string]*Label{
					"test": {
						Name:  "foo",
						Value: "bar",
					},
					"baz": {
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

	node, err := models.Serialize(tag, dbf)

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

	h := "5ea9aca60f50b022efed898ed1dac00278d98de4a5b3178d02a469be218486d6"
	if hex.EncodeToString(classEdge.To.Hash) != h {
		t.Fatalf("invalid hash, expected '%s', got '%s'", h, hex.EncodeToString(classEdge.To.Hash))
	}

	// we store the node a first time
	if err := node.SaveTree(); err != nil {
		t.Fatalf("cannot store node")
	}

	// we modify the second attribute
	tag.Attributes[1].Name = "classes"

	newNode, err := models.Serialize(tag, dbf)

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
	nodes, err := orm.Objects[models.Node](dbf, map[string]any{})

	if err != nil {
		t.Fatal(err)
	}

	// we expect 7 nodes as the initial tree has 5 unique ones and the new tree replaces
	// 2 nodes with modified ones...
	if len(nodes) != 7 {
		t.Fatalf("expected 7 nodes, got %d", len(nodes))
	}

	// we check the number of nodes in the database
	edges, err := orm.Objects[models.Edge](dbf, map[string]any{})

	if err != nil {
		t.Fatal(err)
	}

	if len(edges) != 10 {
		t.Fatalf("expected 10 edges, got %d", len(edges))
	}

	// we restore the node from the Graph DB by its ID
	restoredNode, err := models.GetGraphByID(dbf, node.ID)

	if err != nil {
		t.Fatal(err)
	}

	if restoredNode.ID != node.ID {
		t.Fatalf("IDs do not match")
	}

	// we serialize the original node to JSON
	nodeData, err := json.MarshalIndent(node, " ", " ")

	if err != nil {
		t.Fatal(err)
	}

	// we serialize the restored node to JSON
	restoredNodeData, err := json.MarshalIndent(restoredNode, " ", " ")

	if err != nil {
		t.Fatal(err)
	}

	// we compare the data of the restored node with the original one
	if string(nodeData) != string(restoredNodeData) {
		t.Fatalf("restored node does not match:\n%s\n----\n%s", string(nodeData), string(restoredNodeData))
	}

}
