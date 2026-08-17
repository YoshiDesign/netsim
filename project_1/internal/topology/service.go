package topology

import (
	"errors"
	"fmt"
	"net/netip"
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

	// Status codes are provided in comments as a reference

	// Topo errors
	ErrTopologyNotFound     = errors.New("topology not found")        // use when there's no topology for a given id - 404
	ErrTopologyNameRequired = errors.New("topology name is required") // 400
	ErrTopologyIdNotFound   = errors.New("topology ID not found")     // Use when there's not topo-ID in the request - 400

	// Node errors
	ErrNodeNotFound    = errors.New("node not found")            // 404
	ErrEmptyNodeName   = errors.New("node name cannot be empty") // 400
	ErrInvalidNodeType = errors.New("invalid node type")         // 400
	ErrDuplicateNode   = errors.New("node name already exists")  // 400

	// Link errors
	ErrEmptyLinkEndpoint = errors.New("link endpoint cannot be empty") // 400
	ErrSelfLink          = errors.New("node cannot link to itself")    // 400
	ErrDuplicateLink     = errors.New("link already exists")           // 400 (409)
	ErrLinkNotFound      = errors.New("link not found")                // 404

	// Interface errors
	ErrEmptyInterfaceName     = errors.New("interface name cannot be empty") // 400
	ErrDuplicateInterfaceName = errors.New("interface name already exists")  // 400 (409)
	ErrInterfaceNotFound      = errors.New("interface not found")            // 404

	// IP
	ErrInvalidIPv4Prefix = errors.New("invalid ipv4 address prefix")
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
func (s *TopologyService) NextNodeId() NodeID {
	id := s.nextNodeNum.Add(1)
	return NodeID(fmt.Sprintf("node-%d", id))
}

func (s *TopologyService) NextInterfaceId() InterfaceID {
	id := s.nextIfaceNum.Add(1)
	return InterfaceID(fmt.Sprintf("iface-%d", id))
}

func (s *TopologyService) NextLinkId() LinkID {
	id := s.nextLinkNum.Add(1)
	return LinkID(fmt.Sprintf("link-%d", id))
}

func (s *TopologyService) GetNode(topologyId string, nodeId NodeID) (Node, error) {

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
func (s *TopologyService) CreateNode(topologyId string, name NodeName, nodeType NodeType) (Node, error) {

	if name.TrimSpace() == "" {
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

func (s *TopologyService) DeleteNode(topologyId string, nodeId NodeID) error {

	return s.store.UpdateTopology(
		topologyId,
		func(topo *Topology) error {

			node, ok := topo.Nodes[nodeId]
			if !ok {
				return ErrNodeNotFound
			}

			// Delete links to/from this node
			for _, link := range topo.Links {
				if link.NodeAName == node.Name || link.NodeBName == node.Name {
					delete(topo.Links, link.ID)
				}
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
	nodeA NodeID,
	nodeB NodeID,
) (Link, error) {

	nodeA = nodeA.TrimSpace()
	nodeB = nodeB.TrimSpace()

	if nodeA == "" || nodeB == "" {
		return Link{}, ErrEmptyLinkEndpoint
	}

	if nodeA == nodeB {
		return Link{}, ErrSelfLink
	}

	// Normalize the requested link endpoint for storage
	nodeA, nodeB = canonicalEndpointsByID(nodeA, nodeB)

	link := Link{
		ID:        s.NextLinkId(), // increment our monotonic ID
		NodeAID:   nodeA,
		NodeBID:   nodeB,
		NodeAName: NodeName(INVALID_NODE),
		NodeBName: NodeName(INVALID_NODE),
	}

	// Mutate Topology - Create Links - Critical Section
	// All invariants involving mutable topology state need to be checked inside its callback
	err := s.store.UpdateTopology(
		topoId,
		func(topo *Topology) error {

			// Search for the  nodes in the topology
			foundA := false
			foundB := false
			for _, n := range topo.Nodes {
				if !foundA {
					foundA = n.ID == nodeA
					// Assign the (A) node's ID with the link
					link.NodeAName = n.Name
				}
				if !foundB {
					foundB = n.ID == nodeB
					// Assign the (B) node's ID with the link
					link.NodeBName = n.Name
				}
			}

			// TODO: This should be improved upon - report which was missing
			if !foundA {
				return fmt.Errorf("%w: %s", ErrNodeNotFound, nodeA)
			}
			if !foundB {
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

func (s *TopologyService) GetLink(topologyId string, linkId LinkID) (Link, error) {
	if strings.TrimSpace(topologyId) == "" {
		return Link{}, ErrTopologyIdNotFound
	}

	if linkId.TrimSpace() == "" {
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

func (s *TopologyService) DeleteLink(topologyId string, linkId LinkID) error {
	if strings.TrimSpace(topologyId) == "" {
		return ErrTopologyIdNotFound
	}

	if linkId.TrimSpace() == "" {
		return ErrLinkNotFound
	}

	// Critical Section
	// All invariants involving mutable topology state need to be checked inside its callback
	return s.store.UpdateTopology(
		topologyId,
		func(topo *Topology) error {

			if _, ok := topo.Links[linkId]; !ok {
				return ErrLinkNotFound
			}

			// Note how simple this is with the support
			// of our UpdateTopology abstraction!
			delete(topo.Links, linkId)

			return nil
		})
	// Critical Section
}

/**/
/* * * * * *
* Interfaces
 */
/**/

func (s *TopologyService) CreateInterface(
	topologyId string,
	nodeID NodeID,
	name InterfaceName,
) (Interface, error) {

	if name.TrimSpace() == "" {
		return Interface{}, ErrEmptyInterfaceName
	}

	// Does topology exist?
	//         ↓
	// Does node exist?
	//         ↓
	// Is name valid?
	//         ↓
	// Does node already have an interface with that name?
	//         ↓
	// Create interface

	iface := Interface{
		Name: name,
		ID:   s.NextInterfaceId(),
	}

	err := s.store.AddInterface(topologyId, nodeID, iface)
	if err != nil {
		return Interface{}, err
	}

	// TODO
	// enforce uniqueness
	// construct interface
	// persist mutation

	return iface, nil
}

/* * * * * * *
* TODO - Many of these functions might belong in a separate NetworkingService
* which utilizes the same underlying store
 */
/**/
func (s *TopologyService) SetInterfaceAddress(
	topologyID string,
	nodeID NodeID,
	interfaceID InterfaceID,
	address string,
) (Interface, error) {

	prefix, err := netip.ParsePrefix(address)
	if err != nil {
		// %q preserves any control sequences or escaped values. Great for debugging
		return Interface{}, fmt.Errorf("%w: %q", ErrInvalidIPv4Prefix, address)
	}

	// We're rejecting IPv6 for now
	if !prefix.Addr().Is4() {
		return Interface{}, fmt.Errorf("%w: %q", ErrInvalidIPv4Prefix, address)
	}

	var updated Interface

	// Using our critical section update
	s.store.UpdateTopology(topologyID, func(topo *Topology) error {
		// get the node
		node, ok := topo.Nodes[nodeID]
		if !ok {
			return ErrNodeNotFound
		}

		// Locate the interface we're
		for idx := range node.Interfaces {
			iface := &node.Interfaces[idx] // acquire ref

			if iface.ID != interfaceID {
				continue
			}

			// Set the interface address
			iface.Address = prefix
			updated = *iface // copy

			return nil
		}

		return ErrInterfaceNotFound

	})

	return updated, nil

}
