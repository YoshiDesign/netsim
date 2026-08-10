package store

import (
	"fmt"
	"netsim_0/internal/topology"
	"sync"
)

/*
* Note: Memory store will eventually become an interface
* as we will eventually work with multiple different stores
 */
type MemoryStore struct {
	mu sync.RWMutex

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
		result = append(result, topo)
	}

	return result

}

func (s *MemoryStore) GetTopology(id string) (topology.Topology, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	topo, exists := s.topologies[id]
	// if !exists {
	// 	return topology.Topology{}, true
	// }

	return topo, exists
}

func (s *MemoryStore) CreateTopology(name string) topology.Topology {

	// write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("topology-%d", s.nextID)
	s.nextID++

	topology := topology.Topology{
		ID:   id,
		Name: name,
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
