package store

import (
	"fmt"
	"netsim_1/internal/topology"
	"sync"
)

/*
* The repository layer.
*
* Note: Memory store satisfies a documented interface
* as we will eventually work with multiple different stores
*
* Invariant: Critical access happens while holding the correct lock
 */
type MemoryStore struct {
	mu sync.RWMutex

	// locked by mu^
	topologies map[string]topology.Topology
	nextID     uint64
}

func MakeStore() *MemoryStore {
	return &MemoryStore{
		topologies: make(map[string]topology.Topology),
		nextID:     1,
	}
}

func (s *MemoryStore) ListTopologies() []topology.Topology {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]topology.Topology, 0, len(s.topologies))

	// Copy to output
	for _, topo := range s.topologies {
		result = append(result, topo.Clone())
	}

	return result

}

/**
* Find a Topology - synchronized
 */
func (s *MemoryStore) GetTopology(id string) (topology.Topology, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	topo, exists := s.topologies[id]
	if !exists {
		return topology.Topology{}, topology.ErrTopologyNotFound
	}

	// Clone!! Do not transmit references to underlying data!
	return topo.Clone(), nil
}

/**
* Create a topology - synchronized
 */
func (s *MemoryStore) CreateTopology(name string) topology.Topology {

	// write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("topology-%d", s.nextID)
	s.nextID++

	topology := topology.Topology{
		ID:    id,
		Name:  name,
		Nodes: make(map[topology.NodeID]topology.Node),
		Links: make(map[topology.LinkID]topology.Link),
	}

	s.topologies[id] = topology

	return topology
}

/**
* Delete a topology - synchronized
 */
func (s *MemoryStore) DeleteTopology(id string) bool {

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.topologies[id]; !exists {
		return false
	}
	delete(s.topologies, id)

	return true
}

/**
* Any update to Topology state - synchronized
* Race-free updates to topologies in the store
*
* All invariants involving mutable topology state need to
* be checked inside its callback so they'll also be synchronous
 */
func (s *MemoryStore) UpdateTopology(
	id string,
	fn func(*topology.Topology) error,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	topo, exists := s.topologies[id]
	if !exists {
		return topology.ErrTopologyNotFound
	}

	// Perform the mutation/operation. `fn` MUST ENFORCE YOUR INVARIANTS - CONCURRENCY IS NOT A DRILL
	if err := fn(&topo); err != nil {
		return err
	}

	s.topologies[id] = topo

	return nil
}

/*
* Creating an interface requires us to check invariants
* before modifying nodes. This implies a critical section
* in order to 1. Verify invariants and 2. Create the interface
* as one race-free operation
 */
func (s *MemoryStore) AddInterface(
	topologyID string,
	nodeID topology.NodeID,
	iface topology.Interface,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	// topology lookup
	// node lookup
	// duplicate check
	// append
	// Does topology exist?
	//         ↓
	// Does node exist?
	//         ↓
	// Is name valid?
	//         ↓
	// Does node already have an interface with that name?
	//         ↓
	// Create interface

	// Don't use GetTopology, we already have the lock (will cause a deadlock)
	topo, ok := s.topologies[topologyID]
	if !ok {
		return topology.ErrTopologyNotFound
	}

	node, ok := topo.Nodes[nodeID]
	if !ok {
		return topology.ErrNodeNotFound
	}

	node.Interfaces = append(node.Interfaces, iface)
	topo.Nodes[nodeID] = node
	s.topologies[topologyID] = topo

	return nil
}
