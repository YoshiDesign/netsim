package tests

import (
	"errors"
	"netsim_0/internal/store"
	"netsim_0/internal/topology"
	"testing"
)

func TestCreateTopologyRejectsEmptyName(t *testing.T) {
	store := store.MakeStore()
	service := topology.NewTopologyService(store)

	_, err := service.CreateTopology("   ")

	if !errors.Is(err, topology.ErrTopologyNameRequired) {
		t.Fatalf(
			"expected ErrTopologyNameRequired, got %v",
			err,
		)
	}

	topologies := store.ListTopologies()

	if len(topologies) != 0 {
		t.Fatalf(
			"expected 0 topologies, got %d",
			len(topologies),
		)
	}
}
