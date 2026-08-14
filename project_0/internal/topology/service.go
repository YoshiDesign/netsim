package topology

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

/*
* The service layer enforces the domain rules and policies
* before forwarding requests to the API layer.
*
* Note the use of errors.Is() something I learned here is that
* this works even if returned values are via fmt.Errorf. Pretty cool
 */

// Domain specific (non-HTTP) errors for the service layer
var (
	// Topo errors
	ErrTopologyNotFound     = errors.New("topology not found")
	ErrTopologyNameRequired = errors.New("topology name is required")
	ErrTopologyIdNotFound   = errors.New("topology ID not found")

	// Node errors
	ErrNodeNotFound    = errors.New("node not found")
	ErrEmptyNodeName   = errors.New("node name cannot be empty")
	ErrInvalidNodeType = errors.New("invalid node type")
	ErrDuplicateNode   = errors.New("node name already exists")
)

/* */
func NewTopologyService(store TopologyStore) *TopologyService {
	return &TopologyService{
		store: store,
	}
}

/* * * * * *
* Topologies
 */
func (s *TopologyService) GetTopology(id string) (Topology, error) {
	t, exists := s.store.GetTopology(id)
	if !exists {
		return Topology{}, ErrTopologyNotFound
	}

	return t, nil
}

/* */
func (s *TopologyService) CreateTopology(name string) (Topology, error) {
	// Basic name validation
	name = strings.TrimSpace(name)
	if name == "" {
		return Topology{}, ErrTopologyNameRequired
	}

	created := s.store.CreateTopology(name)

	return created, nil
}

/* */
func (s *TopologyService) ListTopologies() []Topology {
	return s.store.ListTopologies()
}

/* */
func (s *TopologyService) DeleteTopology(id string) error {
	deleted := s.store.DeleteTopology(id)
	if !deleted {
		return ErrTopologyNotFound
	}

	return nil
}

/* * * *
* Nodes
 */

/* Monotonic NodeId for Node Creation - Could also fit into the Store but that's more upon interface tracking. */
func (s *TopologyService) NextNodeId() string {
	id := s.nextNodeNum.Add(1)
	return fmt.Sprintf("node-%d", id)
} // noexcept :)

func (s *TopologyService) GetNode(topologyId string, nodeId string) (Node, error) {

	// Get a clone of the topology
	topo, exists := s.store.GetTopology(topologyId)
	if !exists {
		return Node{}, ErrTopologyNotFound
	}

	// Get the node - no need to copy/clone
	node, ok := topo.Nodes[nodeId]
	if !ok {
		return Node{}, ErrNodeNotFound
	}

	return node, nil
}

/**
* Create a node in a topology. Mutates a topology race-free using UpdateTopology
 */
func (s *TopologyService) CreateNode(topologyId string, name string, nodeType NodeType) (Node, error) {

	if strings.TrimSpace(name) == "" {
		return Node{}, ErrEmptyNodeName
	}

	if !nodeType.Valid() {
		return Node{}, ErrInvalidNodeType
	}

	node := Node{
		ID:   s.NextNodeId(),
		Name: name,
		Type: nodeType,
	}

	// Mutate Topology - Critical Section
	err := s.store.UpdateTopology(
		topologyId,
		func(topo *Topology) error {

			for _, existing := range topo.Nodes {
				// Look for name collision
				if existing.Name == node.Name {
					return ErrDuplicateNode
				}

			}
			topo.Nodes[node.ID] = node
			return nil
		})
	// Critical Section

	if err != nil {
		return Node{}, err
	}

	return node, nil
}

/**
* List all nodes in a Topology -
* Nodes will be sorted for deterministic output, since map access doesn't
* care about order in Go
 */
func (s *TopologyService) ListNodes(topologyId string) ([]Node, error) {
	topo, ok := s.store.GetTopology(topologyId)
	if !ok {
		return nil, ErrTopologyNotFound
	}

	nodes := make([]Node, 0, len(topo.Nodes))

	// Copy each node into output
	for _, node := range topo.Nodes {
		nodes = append(nodes, node)
	}

	// Sort for deterministic/predictable output
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID < nodes[j].ID
	})

	return nodes, nil
}

func (s *TopologyService) DeleteNode(topologyId string, nodeId string) error {

	return s.store.UpdateTopology(
		topologyId,
		func(topo *Topology) error {
			if _, ok := topo.Nodes[nodeId]; !ok {
				return ErrNodeNotFound
			}

			// Note how simple this is with the support
			// of our UpdateTopology abstraction!
			delete(topo.Nodes, nodeId)

			return nil
		})
}
