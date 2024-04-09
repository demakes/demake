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
