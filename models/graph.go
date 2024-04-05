package models

import (
	"github.com/gospel-sh/gospel/orm"
)

type Node struct {
	orm.DBModel
	orm.JSONModel
	Type     string  `json:"type"`
	Hash     []byte  `json:"hash"`
	Outgoing []*Edge `json:"outgoing" db:"ignore"`
	Incoming []*Edge `json:"incoming" db:"ignore"`
}

func MakeNode(db func() orm.DB) *Node {
	return orm.Init(&Node{}, db)
}

type Edge struct {
	orm.DBModel
	orm.JSONModel
	Name   string `json:"name"`
	Key    string `json:"key"`
	Index  string `json:"ind" db:"col:ind"`
	FromID int64  `json:"fromID"`
	From   *Node  `db:"fk:FromID"`
	ToID   int64  `json:"toID"`
	To     *Node  `db:"fk:ToID"`
}

func MakeEdge(db func() orm.DB) *Edge {
	return orm.Init(&Edge{}, db)
}

func (e *Edge) FromTo(from, to *Node) {
	e.From = from
	e.To = from
	e.FromID = from.ID
	e.ToID = to.ID
	from.Outgoing = append(from.Outgoing, e)
	from.Incoming = append(from.Incoming, e)
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
