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

// Domain specific (non-HTTP) invariants for the service layer
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

	// Link errors
	ErrEmptyLinkEndpoint = errors.New("link endpoint cannot be empty")
	ErrSelfLink          = errors.New("node cannot link to itself")
	ErrDuplicateLink     = errors.New("link already exists")
	ErrLinkNotFound      = errors.New("link not found")
)

/* */
func NewTopologyService(store TopologyStore) *TopologyService {
	return &TopologyService{
		store: store,
	}
}

/**/
/* * * * * *
* Topologies
 */
/**/
func (s *TopologyService) GetTopology(id string) (Topology, error) {
	t, err := s.store.GetTopology(id)
	if err != nil {
		return Topology{}, err
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

	// Important that we clone. Go will return a reference
	// to slices & maps from the underlying data
	return created.Clone(), nil
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

/**/
/* * * *
* Nodes
 */
/**/
/* Monotonic NodeId for Node Creation - Could also fit into the Store but that's more upon interface tracking. */
func (s *TopologyService) NextNodeId() string {
	id := s.nextNodeNum.Add(1)
	return fmt.Sprintf("node-%d", id)
} // noexcept :)

func (s *TopologyService) GetNode(topologyId string, nodeId string) (Node, error) {

	// Get a clone of the topology
	topo, err := s.store.GetTopology(topologyId)
	if err != nil {
		return Node{}, err
	}

	// Get the node
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

			// Validation - node's dont already exist by name
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
	topo, err := s.store.GetTopology(topologyId)
	if err != nil {
		return nil, err
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

/**/
/* * * *
* Links
 */
/**/
func (s *TopologyService) CreateLink(
	topoId string,
	nodeA string,
	nodeB string,
) (Link, error) {

	nodeA = strings.TrimSpace(nodeA)
	nodeB = strings.TrimSpace(nodeB)

	if nodeA == "" || nodeB == "" {
		return Link{}, ErrEmptyLinkEndpoint
	}

	if nodeA == nodeB {
		return Link{}, ErrSelfLink
	}

	// Normalize the requested link endpoint for storage
	nodeA, nodeB = canonicalEndpoints(nodeA, nodeB)

	link := Link{
		ID:    fmt.Sprintf("link-%d", s.nextLinkNum.Add(1)), // increment our monotonic ID
		NodeA: nodeA,
		NodeB: nodeB,
	}

	// Mutate Topology - Create Links - Critical Section
	// All invariants involving mutable topology state need to be checked inside its callback
	err := s.store.UpdateTopology(
		topoId,
		func(topo *Topology) error {

			if _, exists := topo.Nodes[nodeA]; !exists {
				return fmt.Errorf("%w: %s", ErrNodeNotFound, nodeA)
			}

			if _, exists := topo.Nodes[nodeB]; !exists {
				return fmt.Errorf("%w: %s", ErrNodeNotFound, nodeB)
			}

			// Validate that the links don't already exist
			if _, exists := topo.Links[link.ID]; exists {
				return ErrDuplicateLink
			}

			topo.Links[link.ID] = link
			return nil
		})
	// Critical Section

	if err != nil {
		return Link{}, err
	}
	return link, nil

}

func (s *TopologyService) ListLinks(topologyId string) ([]Link, error) {

	if strings.TrimSpace(topologyId) == "" {
		return nil, ErrTopologyIdNotFound
	}

	topo, err := s.store.GetTopology(topologyId)
	if err != nil {
		return nil, err
	}

	links := make([]Link, 0, len(topo.Links))

	for _, l := range topo.Links {
		links = append(links, l)
	}

	// Sort for deterministic/predictable output
	sort.Slice(links, func(i, j int) bool {
		return links[i].ID < links[j].ID
	})

	return links, nil

}

func (s *TopologyService) GetLink(topologyId string, linkId string) (Link, error) {
	if strings.TrimSpace(topologyId) == "" {
		return Link{}, ErrTopologyIdNotFound
	}

	if strings.TrimSpace(linkId) == "" {
		return Link{}, ErrLinkNotFound
	}

	topo, err := s.store.GetTopology(topologyId)
	if err != nil {
		return Link{}, err
	}

	// Get the link
	link, exists := topo.Links[linkId]
	if !exists {
		return Link{}, ErrLinkNotFound
	}

	return link, nil
}

func (s *TopologyService) DeleteLink(topologyId string, linkId string) error {
	if strings.TrimSpace(topologyId) == "" {
		return ErrTopologyIdNotFound
	}

	if strings.TrimSpace(linkId) == "" {
		return ErrLinkNotFound
	}

	// Critical Section
	// All invariants involving mutable topology state need to be checked inside its callback
	return s.store.UpdateTopology(
		topologyId,
		func(topo *Topology) error {

			if _, ok := topo.Links[linkId]; !ok {
				return ErrNodeNotFound
			}

			// Note how simple this is with the support
			// of our UpdateTopology abstraction!
			delete(topo.Links, linkId)

			return nil
		})
	// Critical Section
}
