package tests

import (
	"errors"
	"netsim_0/internal/store"
	"netsim_0/internal/topology"
	"testing"
)

// TestCreateLink
// TestCreateLinkTopologyNotFound
// TestCreateLinkNodeANotFound
// TestCreateLinkNodeBNotFound
// TestCreateLinkSelfReference
// TestCreateLinkDuplicate
// TestCreateLinkDuplicateReversed
// TestDeleteLink
// TestDeleteLinkNotFound
// TestDeleteNodeRemovesAttachedLinks

/*
* Here it's important to test the actual graph invariant,
* not just simple Create/Update/Delete logic.
* For this reason it's crucial to involve the service layer
* with each relevant test.
 */

func TestCreateLink(t *testing.T) {

}
func TestCreateLinkTopologyNotFound(t *testing.T) {

}
func TestCreateLinkNodeANotFound(t *testing.T) {

}
func TestCreateLinkNodeBNotFound(t *testing.T) {

}
func TestCreateLinkSelfReference(t *testing.T) {

}
func TestCreateLinkDuplicate(t *testing.T) {

}
func TestCreateLinkDuplicateReversed(t *testing.T) {

	store := store.MakeStore()
	service := topology.NewTopologyService(store)

	testTopology, err := service.CreateTopology("test-topology")
	if err != nil {
		t.Fatal(err)
	}

	topoId := testTopology.ID

	node1, err := service.CreateNode(topoId, "test-node-1", topology.NodeTypeRouter)
	if err != nil {
		t.Fatal(err)
	}
	node2, err := service.CreateNode(topoId, "test-node-2", topology.NodeTypeSwitch)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.CreateLink(
		topoId,
		node1.ID,
		node2.ID,
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.CreateLink(
		topoId,
		node2.ID,
		node1.ID,
	)

	if err != nil && !errors.Is(err, topology.ErrDuplicateLink) {
		t.Fatalf("expected ErrDuplicateLink, got %v", err)
	}
}
func TestDeleteLink(t *testing.T) {

}
func TestDeleteLinkNotFound(t *testing.T) {

}
func TestDeleteNodeRemovesAttachedLinks(t *testing.T) {

}
