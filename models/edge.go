package models

import (
	"encoding/json"
	"fmt"
	"github.com/gospel-sh/gospel/orm"
	"time"
)

type Edge struct {
	orm.DBBaseModel
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
