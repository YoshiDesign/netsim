package tests

import (
	"fmt"
	"netsim_0/internal/store"
	"netsim_0/internal/topology"
	"sync"
	"testing"
)

func TestConcurrentNodeCreation(t *testing.T) {
	// create topology

	const workerCount = 100

	wg := sync.WaitGroup{}

	// Dependencies
	store := store.MakeStore()
	topoService := topology.NewTopologyService(store)
	topo, err := topoService.CreateTopology("test")
	if err != nil {
		t.Fatal(err)
	}

	node, err := topoService.CreateNode(
		topo.ID,
		"router-1",
		topology.NodeTypeRouter,
	)
	if err != nil {
		t.Fatal(err)
	}

	// This is likely the most important test we have so far.
	// Mandatory 100% coverage of mutations.
	if len(topo.Nodes) != 0 {
		t.Fatalf(
			"returned topology was mutated through shared map storage",
		)
	}

	err = topoService.DeleteNode(topo.ID, node.ID)
	if err != nil {
		t.Fatalf(
			"failed to delete node",
		)
	}

	// Create 100 routers concurrently
	for i := range workerCount {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			nodeName := fmt.Sprintf("router-%d", i)

			node, err := topoService.CreateNode(topo.ID, topology.NodeName(nodeName), topology.NodeTypeRouter)
			if err != nil {
				t.Errorf("fail: %v", err)
			}

			if node.Name != topology.NodeName(nodeName) {
				t.Error("fail: node name does not match request")
			}

		}(i)
	}

	wg.Wait()

	updatedTopo, err := topoService.GetTopology(topo.ID)
	if err != nil {
		t.Fatalf("failed to retrieve topology: %v", err)
	}

	if len(updatedTopo.Nodes) != workerCount {
		t.Fatalf(
			"expected %d nodes, got %d",
			workerCount,
			len(updatedTopo.Nodes),
		)
	}

	idMap := make(map[topology.NodeID]bool)
	for _, node := range topo.Nodes {

		if exists := idMap[node.ID]; exists {
			t.Error("fail: duplicate node ID found")
		}

		idMap[node.ID] = true
	}
}
