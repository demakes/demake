package models

import (
	"encoding/json"
	"fmt"
	"github.com/gospel-sh/gospel/orm"
	"time"
)

type Node struct {
	orm.DBBaseModel
	Data     []byte `json:"data" db:"col:data"`
	ID       int64  `json:"id" db:"pk,auto,noOnConflict"`
	Type     string `json:"type"`
	Hash     []byte `json:"hash" db:"pk"`
	Outgoing Edges  `json:"outgoing" db:"ignore"`
	Incoming Edges  `json:"-" db:"ignore"`
}

func (n *Node) SetData(data any) error {
	if bytes, err := json.Marshal(data); err != nil {
		return err
	} else {
		n.Data = bytes
		return nil
	}
}

type Edges []*Edge

func (e Edges) FilterByName(name string) Edges {
	filteredEdges := make(Edges, 0, len(e))

	for _, edge := range e {
		if edge.Name == name {
			filteredEdges = append(filteredEdges, edge)
		}
	}
	return filteredEdges
}

var insertNodeQuery = `
INSERT INTO node
	(
		hash,
		type,
		data,
		updated_at
	)
VALUES
	(
		$1,
		$2,
		$3,
		NULL
	)
ON CONFLICT
	(hash)
WHERE
	deleted_at IS NULL
DO UPDATE SET updated_at = $4
RETURNING
	id, updated_at, created_at
`

func (n *Node) SaveTree(db orm.Transaction) error {

	n.UpdatedAt = &orm.Time{time.Now()}

	if rows, err := db.Query(insertNodeQuery, n.Hash, n.Type, n.Data, n.UpdatedAt.Get()); err != nil {
		return fmt.Errorf("cannot check for node existence. %v", err)
	} else if rows.Next() {
		if err := rows.Scan(&n.ID, &n.UpdatedAt, &n.CreatedAt); err != nil {
			return fmt.Errorf("cannot scan ID: %v", err)
		}
		rows.Close()
		if n.UpdatedAt != nil {
			// the node already exists
			return nil
		}
		// we also need to save all outgoing edges
		for _, edge := range n.Outgoing {
			// we first save the related node and its descendants
			if err := edge.To.SaveTree(db); err != nil {
				return err
			}
			// then we save the edge knowing all nodes exist
			if err := edge.Save(db); err != nil {
				return fmt.Errorf("cannot save edge: %v", err)
			}
		}
	} else {
		rows.Close()
	}
	return nil

}

type Edge struct {
	orm.DBModel
	Data   []byte `json:"data" db:"col:data"`
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Key    string `json:"key"`
	Index  int    `json:"index" db:"col:ind"`
	FromID int64  `json:"fromID"`
	ToID   int64  `json:"toID"`
	Follow bool   `json:"follow"`
	From   *Node  `db:"fk:FromID" json:"-"`
	To     *Node  `db:"fk:ToID" json:"to"`
}

var insertEdgeQuery = `
INSERT INTO edge
	(
		ext_id,
		from_id,
		to_id,
		name,
		type,
		ind,
		key,
		data,
		updated_at
	)
VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		NULL
	)
ON CONFLICT
	(from_id, to_id, name, ind, key, type)
WHERE
	deleted_at IS NULL
DO UPDATE SET updated_at = $9
RETURNING
	id, updated_at, created_at
`

func MakeEdge() *Edge {
	return &Edge{
		Follow: true,
	}
}

func (e *Edge) SetData(data any) error {
	if bytes, err := json.Marshal(data); err != nil {
		return err
	} else {
		e.Data = bytes
		return nil
	}
}

func (e *Edge) Save(db orm.Transaction) error {
	if e.From == nil || e.From.ID == 0 {
		return fmt.Errorf("'From' node missing or doesn't have an ID")
	}
	if e.To == nil || e.To.ID == 0 {
		return fmt.Errorf("'To' node missing or doesn't have an ID")
	}
	// we update 'From' and 'To' IDs
	e.FromID = e.From.ID
	e.ToID = e.To.ID

	if e.ExtID == nil {
		e.ExtID = &orm.UUID{}
		if err := e.ExtID.Generate(); err != nil {
			return err
		}
	}

	e.UpdatedAt = &orm.Time{time.Now()}

	if rows, err := db.Query(insertEdgeQuery, e.ExtID.Bytes(), e.FromID, e.ToID, e.Name, e.Type, e.Index, e.Key, e.Data, time.Now().UTC()); err != nil {
		return fmt.Errorf("cannot check for node existence. %v", err)
	} else {
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&e.ID, &e.UpdatedAt, &e.CreatedAt); err != nil {
				return fmt.Errorf("cannot scan ID: %v", err)
			}
		} else {
			return fmt.Errorf("cannot insert edge")
		}
	}
	return nil
}

func (e *Edge) FromTo(from, to *Node) {
	e.From = from
	e.To = to
	e.FromID = from.ID
	e.ToID = to.ID
	from.Outgoing = append(from.Outgoing, e)
	to.Incoming = append(to.Incoming, e)
}
