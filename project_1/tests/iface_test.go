package tests

import (
	"errors"
	"netsim_1/internal/store"
	"netsim_1/internal/topology"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentDuplicateInterfaceCreation(t *testing.T) {

	store := store.MakeStore()
	service := topology.NewTopologyService(store)

	topo, _ := service.CreateTopology("Test-Topology-Ifaces")
	topologyID := topo.ID

	node, _ := service.CreateNode(topologyID, "test-node", topology.NodeTypeHost)
	nodeID := node.ID

	const workers = 100

	var wg sync.WaitGroup
	wg.Add(workers)

	var successes atomic.Int32

	// Using 100 concurrent workers, attempt to create the same interface.
	// Only 1 should succeed
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			_, err := service.CreateInterface(
				topologyID,
				nodeID,
				"eth0",
			)

			if err == nil {
				successes.Add(1)
				return
			}

			if !errors.Is(err, topology.ErrDuplicateInterfaceName) {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	if got := successes.Load(); got != 1 {
		t.Fatalf("expected exactly 1 successful create, got %d", got)
	}
}
