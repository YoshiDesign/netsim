package store

import (
	"fmt"
	"sync"
	"testing"
)

/*
* Run these with race detection!
* go test -race ./...
 */

func TestMemoryStoreConcurrentCreate(t *testing.T) {
	store := MakeStore()

	const workers = 100

	var wg sync.WaitGroup
	wg.Add(workers)

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
}
