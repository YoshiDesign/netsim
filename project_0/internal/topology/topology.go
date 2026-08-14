package topology

import "sync/atomic"

// Describe what the API accepts
type CreateTopologyRequest struct {
	Name string `json:"name"`
}

type Topology struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Nodes map[string]Node `json:"nodes"`
	// ^ Bear in mind the referential nature of this field.
	// This is why we provide an explicit Clone() operation
}

// TopologyStore is a consumer of a Store interface
type TopologyStore interface {
	CreateTopology(name string) Topology
	ListTopologies() []Topology
	GetTopology(id string) (Topology, bool)
	DeleteTopology(id string) bool
	UpdateTopology(id string, fn func(topo *Topology) error) error
}

// TopologyServices defines the Service pattern
type TopologyService struct {
	store       TopologyStore
	nextNodeNum atomic.Uint64
}

func (t Topology) Clone() Topology {
	clone := Topology{
		ID:    t.ID,
		Name:  t.Name,
		Nodes: make(map[string]Node, len(t.Nodes)),
	}

	for id, node := range t.Nodes {
		clone.Nodes[id] = node
	}

	return clone
}
