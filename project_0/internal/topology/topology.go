package topology

import "sync/atomic"

type Topology struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Nodes map[NodeID]Node `json:"nodes"`
	Links map[string]Link `json:"links"`
	// ^ Bear in mind the referential nature of maps (& slices).
	// This is why we provide an explicit Clone() operation.
	// Do not return a Topology (or anything containing maps or slices)
	// to the client (and many other callsites) without cloning it first.
}

// TopologyStore is a consumer of a Store interface
type TopologyStore interface {
	CreateTopology(name string) Topology
	ListTopologies() []Topology
	GetTopology(id string) (Topology, error)
	DeleteTopology(id string) bool
	UpdateTopology(id string, fn func(topo *Topology) error) error
}

// TopologyServices defines the Service pattern
type TopologyService struct {
	store       TopologyStore
	nextNodeNum atomic.Uint64
	nextLinkNum atomic.Uint64
}

func (t Topology) Clone() Topology {
	clone := Topology{
		ID:    t.ID,
		Name:  t.Name,
		Nodes: make(map[NodeID]Node, len(t.Nodes)),
		Links: make(map[string]Link, len(t.Links)),
	}

	for id, node := range t.Nodes {
		clone.Nodes[id] = node
	}

	for id, link := range t.Links {
		clone.Links[id] = link
	}

	return clone
}
