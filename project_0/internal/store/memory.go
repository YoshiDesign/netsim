package store

import (
	"fmt"
	"netsim_0/internal/topology"
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
* Find a Topology -
* Note how we return a bool instead of an error. This is more idiomatic
* as this function doesn't provide any more than a found/not-found outcome.
* There's only 1 failure mode.
 */
func (s *MemoryStore) GetTopology(id string) (topology.Topology, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	topo, exists := s.topologies[id]
	if !exists {
		return topology.Topology{}, false
	}

	// Clone!! Do not transmit references to underlying data!
	return topo.Clone(), true
}

func (s *MemoryStore) CreateTopology(name string) topology.Topology {

	// write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("topology-%d", s.nextID)
	s.nextID++

	topology := topology.Topology{
		ID:    id,
		Name:  name,
		Nodes: make(map[string]topology.Node),
	}

	s.topologies[id] = topology

	return topology
}

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
* Race-free updates to topologies in the store
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

	// Perform the mutation/operation
	if err := fn(&topo); err != nil {
		return err
	}

	s.topologies[id] = topo

	return nil
}
