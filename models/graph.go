package models

import (
	"fmt"
	"github.com/gospel-sh/gospel/orm"
)

type Node struct {
	orm.DBBaseModel
	orm.JSONModel
	ID       int64   `json:"-" db:"pk,auto,noOnConflict"`
	Type     string  `json:"type"`
	Hash     []byte  `json:"hash" db:"pk"`
	Outgoing []*Edge `json:"outgoing" db:"ignore"`
	Incoming []*Edge `json:"incoming" db:"ignore"`
}

func MakeNode(db func() orm.DB) *Node {
	return orm.Init(&Node{}, db)
}

func (n *Node) Save() error {
	return orm.Save(n)
}

func (n *Node) Refresh() (bool, error) {
	if err := orm.LoadOne(n, map[string]any{"hash": n.Hash}); err == nil {
		return true, nil
	} else if err == orm.NotFound {
		return false, nil
	} else {
		return false, err
	}
}

func (n *Node) SaveTree() error {
	if ok, err := n.Refresh(); err != nil {
		return fmt.Errorf("cannot check for node existence")
	} else if !ok {
		// this node doesn't exist yet, we save it
		if err := n.Save(); err != nil {
			return fmt.Errorf("cannot save node: %v", err)
		}
		// we also need to save all outgoing edges
		for _, edge := range n.Outgoing {
			// we first save the related node and its descendants
			if err := edge.To.SaveTree(); err != nil {
				return err
			}
			// then we save the edge knowing all nodes exist
			if err := edge.Save(); err != nil {
				return fmt.Errorf("cannot save edge: %v", err)
			}
		}
	}
	return nil

}

type Edge struct {
	orm.DBModel
	orm.JSONModel
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Key    string `json:"key"`
	Index  int    `json:"index" db:"col:ind"`
	FromID int64  `json:"fromID"`
	From   *Node  `db:"fk:FromID" json:"from"`
	ToID   int64  `json:"toID"`
	To     *Node  `db:"fk:ToID" json:"to"`
}

func MakeEdge(db func() orm.DB) *Edge {
	return orm.Init(&Edge{}, db)
}

func (e *Edge) Save() error {
	if e.From == nil || e.To.ID == 0 {
		return fmt.Errorf("'From' node missing or doesn't have an ID")
	}
	if e.To == nil || e.To.ID == 0 {
		return fmt.Errorf("'To' node missing or doesn't have an ID")
	}
	// we update 'From' and 'To' IDs
	e.FromID = e.From.ID
	e.ToID = e.To.ID
	return orm.Save(e)
}

func (e *Edge) FromTo(from, to *Node) {
	e.From = from
	e.To = to
	e.FromID = from.ID
	e.ToID = to.ID
	from.Outgoing = append(from.Outgoing, e)
	to.Incoming = append(to.Incoming, e)
}

type NodeWithEdge struct {
	Edge *Edge
	Node *Node
}

// return the entire graph for a given node
var query = `
WITH RECURSIVE
	graph(name, key, ind, type, data, hash, from_id, to_id)
	AS (
		SELECT  '', '', 0, type, data, hash, 0, id from node where id = 6
		UNION ALL SELECT
			edge.name, edge.key, edge.ind, node.type, node.data, node.hash, edge.from_id, edge.to_id FROM edge JOIN node ON node.id = edge.to_id JOIN graph ON edge.from_id = graph.to_id) select * from graph;
`
