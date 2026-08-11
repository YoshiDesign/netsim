package tests

import (
	"fmt"
	"netsim_0/internal/store"
	"sync"
	"testing"
)

/*
* Run these with race detection!
* go test -race ./...
 */

func TestMemoryStoreConcurrentCreate(t *testing.T) {
	store := store.MakeStore()

	const workers = 100

	var wg sync.WaitGroup
	wg.Add(workers)

	// Part 1 - race check
	for i := range workers {
		go func(id int) {
			defer wg.Done()

			store.CreateTopology(fmt.Sprintf("test-%d", i))

		}(i)
	}

	wg.Wait()

	topologies := store.ListTopologies()

	if len(topologies) != workers {
		t.Fatalf(
			"expected %d topologies, got %d",
			workers,
			len(topologies),
		)
	}

	// Part 2 - make sure every id is unique
	ids := make(map[string]struct{}) // a common Go "set-like" idioms. struct{} consumes no meaningful storage. Points to the zero sentinel

	for _, topology := range topologies {
		if _, exists := ids[topology.ID]; exists {
			t.Fatalf("duplicate topology ID: %s", topology.ID)
		}

		ids[topology.ID] = struct{}{}
	}
}
