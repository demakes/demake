package models_test

import (
	"fmt"
	"github.com/gospel-sh/gospel/orm"
	"github.com/klaro-org/sites"
	"github.com/klaro-org/sites/models"
	kt "github.com/klaro-org/sites/testing"
	"testing"
)

func BenchmarkSimpleSave(b *testing.B) {

	if err := registerModels(); err != nil {
		b.Fatal(err)
	}

	settings, err := sites.LoadSettings()

	if err != nil {
		b.Fatal(err)
	}

	labels := map[string]*Label{}

	for i := 0; i < 100; i++ {
		labels[fmt.Sprintf("%d", i)] = &Label{
			Name:  fmt.Sprintf("foo%d", i),
			Value: fmt.Sprintf("bar%d", i),
		}
	}

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
				Name:   "class",
				Value:  "bar",
				Labels: labels,
			},
		},
	}

	b.ResetTimer()
	b.StopTimer()

	for i := 0; i < b.N; i++ {

		db, err := kt.DB(settings)

		if err != nil {
			b.Fatal(err)
		}

		b.StartTimer()

		tx, err := db.Begin()

		if err != nil {
			b.Fatal(err)
		}

		node, err := models.Serialize(tag)

		if err != nil {
			b.Fatal(err)
		}

		// we store the node a first time
		if err := node.SaveTree(db); err != nil {
			b.Fatalf("cannot store node")
		}

		if err := tx.Commit(); err != nil {
			b.Fatal(err)
		}

		b.StopTimer()

	}

}

func BenchmarkDeepTree(b *testing.B) {

	if err := registerModels(); err != nil {
		b.Fatal(err)
	}

	settings, err := sites.LoadSettings()

	if err != nil {
		b.Fatal(err)
	}

	tag := &Tag{
		Type:       "p",
		Meta:       Meta{Language: "de"},
		Attributes: []*Attribute{},
		Children:   []*Tag{},
	}

	currentChild := tag

	// we create a tree with a depth of 51 nodes
	for i := 0; i < 50; i++ {
		childTag := &Tag{
			Type:       fmt.Sprintf("h%d", i),
			Attributes: []*Attribute{},
			Children:   []*Tag{},
			Meta:       Meta{Language: "de"},
		}
		currentChild.Children = append(currentChild.Children, childTag)
		currentChild = childTag
	}

	b.ResetTimer()
	b.StopTimer()

	for i := 0; i < b.N; i++ {

		db, err := kt.DB(settings)

		if err != nil {
			b.Fatal(err)
		}

		b.StartTimer()

		tx, err := db.Begin()

		if err != nil {
			b.Fatal(err)
		}

		node, err := models.Serialize(tag)

		if err != nil {
			b.Fatal(err)
		}

		// we store the node a first time
		if err := node.SaveTree(db); err != nil {
			b.Fatalf("cannot store node")
		}

		// we modify the innermost child
		currentChild.Type = "foo"

		newNode, err := models.Serialize(tag)

		if err != nil {
			b.Fatal(err)
		}

		// we store the node a first time
		if err := newNode.SaveTree(db); err != nil {
			b.Fatalf("cannot store node")
		}

		if err := tx.Commit(); err != nil {
			b.Fatal(err)
		}

		b.StopTimer()

	}

}

func BenchmarkDeepRead(b *testing.B) {

	if err := registerModels(); err != nil {
		b.Fatal(err)
	}

	settings, err := sites.LoadSettings()

	if err != nil {
		b.Fatal(err)
	}

	tag := &Tag{
		Type:       "p",
		Meta:       Meta{Language: "de"},
		Attributes: []*Attribute{},
		Children:   []*Tag{},
	}

	currentChild := tag

	// we create a tree with a depth of 51 nodes
	for i := 0; i < 100; i++ {
		childTag := &Tag{
			Type:       fmt.Sprintf("h%d", i),
			Attributes: []*Attribute{},
			Children:   []*Tag{},
			Meta:       Meta{Language: "de"},
		}
		currentChild.Children = append(currentChild.Children, childTag)
		currentChild = childTag
	}

	db, err := kt.DB(settings)

	if err != nil {
		b.Fatal(err)
	}

	dbf := func() orm.DB { return db }

	node, err := models.Serialize(tag)

	if err != nil {
		b.Fatal(err)
	}

	// we store the node a first time
	if err := node.SaveTree(db); err != nil {
		b.Fatalf("cannot store node")
	}

	b.ResetTimer()
	b.StopTimer()

	for i := 0; i < b.N; i++ {

		b.StartTimer()

		// we restore the node from the Graph DB by its ID
		_, err := models.GetGraphByID(dbf, node.ID)

		if err != nil {
			b.Fatal(err)
		}

		b.StopTimer()

	}

}

func BenchmarkDeepAndWideRead(b *testing.B) {

	if err := registerModels(); err != nil {
		b.Fatal(err)
	}

	settings, err := sites.LoadSettings()

	if err != nil {
		b.Fatal(err)
	}

	root := &Tag{
		Type:       "root",
		Meta:       Meta{Language: "de"},
		Attributes: []*Attribute{},
		Children:   []*Tag{},
	}

	for i := 0; i < 100; i++ {
		tag := &Tag{
			Type:       "p",
			Meta:       Meta{Language: "de"},
			Attributes: []*Attribute{},
			Children:   []*Tag{},
		}

		currentChild := tag

		// we create a deep tree
		for j := 0; j < 200; j++ {
			childTag := &Tag{
				Type:       fmt.Sprintf("h%d", i),
				Attributes: []*Attribute{},
				Children:   []*Tag{},
				Meta:       Meta{Language: "de"},
			}
			currentChild.Children = append(currentChild.Children, childTag)
			currentChild = childTag
		}

		root.Children = append(root.Children, tag)
	}

	db, err := kt.DB(settings)

	if err != nil {
		b.Fatal(err)
	}

	dbf := func() orm.DB { return db }

	node, err := models.Serialize(root)

	if err != nil {
		b.Fatal(err)
	}

	tx, err := db.Begin()

	if err != nil {
		b.Fatal(err)
	}

	// we store the node a first time
	if err := node.SaveTree(tx); err != nil {
		b.Fatalf("cannot store node")
	}

	if err := tx.Commit(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.StopTimer()

	fmt.Println("Retrieving...")

	for i := 0; i < b.N; i++ {

		b.StartTimer()

		// we restore the node from the Graph DB by its ID
		_, err := models.GetGraphByID(dbf, node.ID)

		if err != nil {
			b.Fatal(err)
		}

		b.StopTimer()

	}

}
