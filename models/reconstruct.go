package models

import (
	"fmt"
	"github.com/gospel-sh/gospel/orm"
)

// Data structure to describe the combined edge and node data
type GraphData struct {
	db func() orm.DB
	// the order of these fields MUST match the order in the CTE below!
	Name          string
	Key           string
	Index         int `db:"col:ind"`
	Type          string
	Data          []byte
	Hash          []byte
	NodeCreatedAt *orm.Time
	NodeUpdatedAt *orm.Time
	FromID        int64
	ToID          int64
	EdgeID        int64
	EdgeType      int
	EdgeExtID     *orm.UUID
	EdgeData      []byte
	EdgeCreatedAt *orm.Time
	EdgeUpdatedAt *orm.Time
	EdgeFollow    bool
}

// returns the entire graph for a given node stopping at non-follow edges
var graphQuery = `
WITH RECURSIVE
	graph( -- all rows here must appear below as well
		name,
		key,
		ind,
		type,
		data,
		hash,
		node_created_at,
		node_updated_at,
		from_id,
		to_id,
		edge_id,
		edge_type,
		edge_ext_id,
		edge_data,
		edge_created_at,
		edge_updated_at,
		edge_follow)
	AS (
		SELECT
			'',
			'',
			0,
			type,
			data,
			hash,
			created_at,
			updated_at,
			0,
			id,
			0,
			0,
			NULL,
			NULL,
			current_timestamp,
			current_timestamp,
			true
		FROM
			node
		WHERE
			id = ?
		UNION ALL SELECT
			edge.name,
			edge.key,
			edge.ind,
			node.type,
			node.data,
			node.hash,
			node.created_at,
			node.updated_at,
			edge.from_id,
			edge.to_id,
			edge.id,
			edge.type,
			edge.ext_id,
			edge.data,
			edge.created_at,
			edge.updated_at,
			edge.follow
		FROM
			edge
		JOIN
			node ON node.id = edge.to_id AND node.deleted_at IS NULL
		JOIN
			graph ON edge.from_id = graph.to_id
		WHERE
			graph.edge_follow = true AND
			edge.deleted_at IS NULL
		ORDER BY
			edge.to_id -- we sort in a depth-first way
	)
SELECT * FROM graph;
`

func (g *GraphData) Database() func() orm.DB {
	return g.db
}

func (g *GraphData) SetDatabase(db func() orm.DB) {
	g.db = db
}

func (g *GraphData) TableName() string {
	return ""
}

func (g *GraphData) SetTableName(name string) {

}

func (g *GraphData) Init() error {
	return nil
}

func reconstructNode(db func() orm.DB, data []*GraphData) (*Node, []*GraphData, error) {

	// we get the first data element and remove it from the list
	nodeData := data[0]
	data = data[1:]

	// we initialize a new node
	node := MakeNode(db)
	node.ID = nodeData.ToID
	node.Hash = nodeData.Hash
	node.Type = nodeData.Type
	node.CreatedAt = nodeData.NodeCreatedAt
	node.UpdatedAt = nodeData.NodeUpdatedAt
	if err := node.JSON.Set(nodeData.Data); err != nil {
		return nil, nil, fmt.Errorf("cannot set JSON: %v", err)
	}

	// this assumes edges are ordered depth-first!!!
	for len(data) > 0 {

		if data[0].FromID != node.ID {
			// the next edge doesn't belong to the same node, we break
			break
		}

		edgeData := data[0]

		// we initialize a new edge
		edge := MakeEdge(db)
		edge.Type = edgeData.EdgeType
		edge.ExtID = edgeData.EdgeExtID
		edge.Name = edgeData.Name
		edge.Key = edgeData.Key
		edge.Index = edgeData.Index
		edge.CreatedAt = edgeData.EdgeCreatedAt
		edge.UpdatedAt = edgeData.EdgeUpdatedAt
		if err := edge.JSON.Set(edgeData.EdgeData); err != nil {
			return nil, nil, fmt.Errorf("cannot set JSON: %v", err)
		}

		var toNode *Node
		var err error

		// we generate a new node and let it consume all its children
		if toNode, data, err = reconstructNode(db, data); err != nil {
			return nil, nil, err
		}

		// we link the edge to the nodes
		edge.FromTo(node, toNode)
	}

	// we return the node and remaining data
	return node, data, nil

}

func GetGraphByID(db func() orm.DB, id int64) (*Node, error) {

	graphDataList, err := orm.Query[GraphData](db, graphQuery, id)

	if err != nil {
		return nil, err
	}

	if len(graphDataList) == 0 {
		return nil, fmt.Errorf("not found")
	}

	node, _, err := reconstructNode(db, graphDataList)

	return node, err
}
