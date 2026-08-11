package topology

// Describe what the API accepts
type CreateTopologyRequest struct {
	Name string `json:"name"`
}

type Topology struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TopologyStore is a consumer of a *Store's interface
type TopologyStore interface {
	CreateTopology(name string) Topology
	ListTopologies() []Topology
	GetTopology(id string) (Topology, bool)
	DeleteTopology(id string) bool
}

type TopologyService struct {
	store TopologyStore
}
