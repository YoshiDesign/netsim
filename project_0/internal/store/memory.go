package store

import (
	"netsim_0/internal/topology"
	"sync"
)

type MemoryStore struct {
	mu sync.RWMutex

	topologies map[string]topology.Topology
	nextID     uint64
}

func MakeStore() MemoryStore {
	return MemoryStore{}
}
