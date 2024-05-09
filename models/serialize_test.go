package models_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/demakes/demake"
	"github.com/demakes/demake/models"
	kt "github.com/demakes/demake/testing"
	"github.com/gospel-sh/gospel/orm"
	"testing"
)

type Plugin interface {
	Editor() any
}

type Site struct {
	Plugins []Plugin `json:"plugins" graph:"include"`
	Name    string   `json:"name"`
	Domain  string   `json:"domain"`
	Root    *Tag     `json:"root"`
}

type Tag struct {
	Type       string       `json:"type"`
	Meta       Meta         `json:"meta"`
	Value      any          `json:"value"`
	Children   []*Tag       `json:"children"`
	Attributes []*Attribute `json:"attributes"`
}

type Meta struct {
	Language string `json:"language"`
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

func registerModels() error {

	// we first register the label model
	if err := models.Register[Label]("label"); err != nil {
		return err
	}

	// then we register the attribute model
	if err := models.Register[Attribute]("attribute"); err != nil {
		return err
	}

	// then we register the meta model
	if err := models.Register[Meta]("meta"); err != nil {
		return err
	}

	// then we register the tag model
	if err := models.Register[Tag]("tag"); err != nil {
		return err
	}

	if err := models.Register[Site]("site"); err != nil {
		return err
	}

	if err := models.Register[RoutesPlugin]("routesPlugin"); err != nil {
		return err
	}

	if err := models.Register[BlogPlugin]("blogPlugin"); err != nil {
		return err
	}

	return nil
}

type RoutesPlugin struct {
	Prefix string `json:"prefix"`
}

type BlogPlugin struct {
	Title string `json:"title"`
}

func (b *BlogPlugin) Editor() any {
	return "bar"
}

func (r *RoutesPlugin) Editor() any {
	return "foo"
}

func TestSite(t *testing.T) {

	if err := registerModels(); err != nil {
		t.Fatal(err)
	}

	site := &Site{
		Plugins: []Plugin{&RoutesPlugin{Prefix: "/test"}, &BlogPlugin{Title: "My fancy blog"}},
	}

	node, err := models.Serialize(site)

	if err != nil {
		t.Fatal(err)
	}

	restoredSite, err := models.DeserializeType[Site](node)

	if err != nil {
		t.Fatal(err)
	}

	if len(restoredSite.Plugins) != 2 {
		t.Fatalf("expected one plugin")
	}

	if restoredSite.Plugins[0].(*RoutesPlugin).Prefix != "/test" {
		t.Fatalf("prefix doesn't match")
	}

	if restoredSite.Plugins[1].(*BlogPlugin).Title != "My fancy blog" {
		t.Fatalf("title doesn't match")
	}

}

func TestSerialize(t *testing.T) {

	if err := registerModels(); err != nil {
		t.Fatal(err)
	}

	settings, err := sites.LoadSettings()

	if err != nil {
		t.Fatal(err)
	}

	db, err := kt.DB(settings)

	if err != nil {
		t.Fatal(err)
	}

	dbf := func() orm.DB { return db }

	tag := &Tag{
		Type: "p",
		Meta: Meta{Language: "de"},
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
						Name:  "fam",
						Value: "bar",
					},
					"baz": {
						Name:  "flip",
						Value: "flop",
					},
				},
			},
		},
	}

	tagSchema := models.SchemaFor(tag)

	if tagSchema == nil {
		t.Fatalf("expected a schema")
	}

	if len(tagSchema.Fields) != 2 {
		t.Fatalf("expected two regular fields")
	}

	node, err := models.Serialize(tag)

	if err != nil {
		t.Fatal(err)
	}

	expected := `{"type":"p","value":null}`

	if string(node.Data) != expected {
		t.Fatalf("data doesn't match: %s vs. %s", string(node.Data), string(expected))
	}

	if len(node.Outgoing) != 3 {
		t.Fatalf("expected 3 outgoing edges, got %d", len(node.Outgoing))
	}

	styleEdge := node.Outgoing[1]
	classEdge := node.Outgoing[2]

	if styleEdge.Index != 0 {
		t.Fatalf("expected 0 index")
	}

	if classEdge.Index != 1 {
		t.Fatalf("expected 1 index, got %d", classEdge.Index)
	}

	expected = `{"name":"style","value":"font-size:12px"}`

	if string(styleEdge.To.Data) != expected {
		t.Fatalf("data doesn't match: %s vs. %s", styleEdge.To.Data, expected)
	}

	expected = `{"name":"class","value":"bar"}`
	if string(classEdge.To.Data) != expected {
		t.Fatalf("data doesn't match: %s vs. %s", string(styleEdge.To.Data), expected)
	}

	h := "faf218cb19d48110fe8be0ac889d9d56056db79803d0f769ad68776d233e0901"
	if hex.EncodeToString(classEdge.To.Hash) != h {
		t.Fatalf("invalid hash, expected '%s', got '%s'", h, hex.EncodeToString(classEdge.To.Hash))
	}

	// we store the node a first time
	if err := node.SaveTree(db); err != nil {
		t.Fatalf("cannot store node: %v", err)
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

	// we deserialize the restored node into a tag
	restoredTag, err := models.DeserializeType[Tag](restoredNode)

	if err != nil {
		t.Fatal(err)
	}

	// we serialize the original tag to JSON
	tagData, err := json.MarshalIndent(tag, " ", " ")

	if err != nil {
		t.Fatal(err)
	}

	// we serialize the restored tag to JSON
	restoredTagData, err := json.MarshalIndent(restoredTag, " ", " ")

	if err != nil {
		t.Fatal(err)
	}

	// we compare the data of the restored tag with the original one
	if string(tagData) != string(restoredTagData) {
		t.Fatalf("restored node does not match:\n%s\n----\n%s", string(tagData), string(restoredTagData))
	}

	// we modify the second attribute
	tag.Attributes[1].Name = "classes"

	newNode, err := models.Serialize(tag)

	if err != nil {
		t.Fatal(err)
	}

	// we revert the change
	tag.Attributes[1].Name = "class"

	// we store the node a second time
	if err := newNode.SaveTree(db); err != nil {
		t.Fatalf("cannot store new node a second time")
	}

	if newNode.Outgoing[1].To.ID != node.Outgoing[1].To.ID {
		t.Fatalf("expected IDs of first attribute to match")
	}

	if newNode.Outgoing[2].To.ID == node.Outgoing[2].To.ID {
		t.Fatalf("expected IDs of second attribute to diverge")
	}

	// we check the number of nodes in the database
	nodes, err := orm.Objects[models.Node](dbf, map[string]any{})

	if err != nil {
		t.Fatal(err)
	}

	// we expect 10 nodes as the initial tree has 8 unique ones and the new tree replaces
	// 2 nodes with modified ones...
	if len(nodes) != 10 {
		t.Fatalf("expected 10 nodes, got %d", len(nodes))
	}

	// we check the number of nodes in the database
	edges, err := orm.Objects[models.Edge](dbf, map[string]any{})

	if err != nil {
		t.Fatal(err)
	}

	if len(edges) != 12 {
		t.Fatalf("expected 12 edges, got %d", len(edges))
	}

}
