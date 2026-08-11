package topology

import "errors"

// Domain specific (non-HTTP) errors
var (
	ErrTopologyNotFound     = errors.New("topology not found")
	ErrTopologyNameRequired = errors.New("topology name is required")
)

// Describe what the API accepts
type CreateTopologyRequest struct {
	Name string `json:"name"`
}

// Describe a topology we already have
type Topology struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TopologyStore is a consumer of MemoryStore's signatures/interface
type TopologyStore interface {
	CreateTopology(name string) Topology
	ListTopologies() []Topology
	GetTopology(id string) (Topology, bool)
	DeleteTopology(id string) bool
}

type TopologyService struct {
	store TopologyStore
}

func NewTopologyService(store TopologyStore) *TopologyService {
	return &TopologyService{
		store: store,
	}
}
