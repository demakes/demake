package models_test

import (
	"fmt"
	"github.com/demakes/demake"
	"github.com/demakes/demake/models"
	kt "github.com/demakes/demake/testing"
	"github.com/gospel-sh/gospel/orm"
	"testing"
)

func TestBasicGraph(t *testing.T) {

	settings, err := sites.LoadSettings()

	if err != nil {
		t.Fatal(err)
	}

	db, err := kt.DB(settings)

	if err != nil {
		t.Fatal(err)
	}

	var previousNode *models.Node

	for i := 0; i < 5; i++ {

		node := &models.Node{
			Type: "test",
			Hash: []byte(fmt.Sprintf("abc%d", i)),
		}

		orm.Init(node, func() orm.DB { return db })

		if err := node.SetData(map[string]any{"foo": "bar"}); err != nil {
			t.Fatal(err)
		}

		if err := orm.Save(node); err != nil {
			t.Fatal(err)
		}

		if previousNode != nil {
			edge := &models.Edge{
				Name:   "test",
				FromID: previousNode.ID,
				ToID:   node.ID,
				Key:    "foo",
			}

			orm.Init(edge, func() orm.DB { return db })

			if err := orm.Save(edge); err != nil {
				t.Fatal(err)
			}

		}

		previousNode = node

	}

}
