package tests

import (
	"errors"
	"netsim_1/internal/store"
	"netsim_1/internal/topology"
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
}
func TestCreateLinkTopologyNotFound(t *testing.T) {
	store := store.MakeStore()
	service := topology.NewTopologyService(store)

	testTopology, err := service.CreateTopology("test-topology")
	if err != nil {
		t.Fatal(err)
	}

	topoId := testTopology.ID

	node1, err := service.CreateNode(topoId, "test-node-1", topology.NodeTypeRouter)
	if err != nil && !errors.Is(err, topology.ErrTopologyNotFound) {
		t.Fatal(err)
	}

	_, err = service.CreateLink(
		topoId,
		"non-existent",
		node1.ID,
	)

	if err != nil && !errors.Is(err, topology.ErrNodeNotFound) {
		t.Fatal(err)
	}

	if err == nil {
		t.Fatal(err)
	}
}
func TestCreateLinkNodeANotFound(t *testing.T) {

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

	_, err = service.CreateLink(
		topoId,
		"non-existent",
		node1.ID,
	)

	if err != nil && !errors.Is(err, topology.ErrNodeNotFound) {
		t.Fatal(err)
	}

	if err == nil {
		t.Fatal(err)
	}
}
func TestCreateLinkNodeBNotFound(t *testing.T) {

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

	_, err = service.CreateLink(
		topoId,
		node1.ID,
		"non-existent",
	)
	if err == nil {
		t.Fatal(err)
	}

}
func TestCreateLinkSelfReference(t *testing.T) {
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

	_, err = service.CreateLink(
		topoId,
		node1.ID,
		node1.ID,
	)
	if err != nil && !errors.Is(err, topology.ErrSelfLink) {
		t.Fatal(err)
	}
	if err == nil {
		t.Fatal(err)
	}
}
func TestCreateLinkDuplicate(t *testing.T) {

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
		node1.ID,
		node2.ID,
	)

	if err != nil && !errors.Is(err, topology.ErrDuplicateLink) {
		t.Fatalf("expected ErrDuplicateLink, got %v", err)
	}
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

	link, err := service.CreateLink(
		topoId,
		node1.ID,
		node2.ID,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = service.DeleteLink(topoId, link.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func TestDeleteLinkNotFound(t *testing.T) {
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

	link, err := service.CreateLink(
		topoId,
		node1.ID,
		node2.ID,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = service.DeleteLink(topoId, link.ID)
	if err != nil {
		t.Fatal(err)
	}

	err = service.DeleteLink(topoId, "non-existent-link")
	if err != nil && !errors.Is(err, topology.ErrLinkNotFound) {
		t.Fatal(err)
	}

	err = service.DeleteLink(topoId, "      ")
	if err != nil && !errors.Is(err, topology.ErrLinkNotFound) {
		t.Fatal(err)
	}

}
func TestDeleteNodeRemovesAttachedLinks(t *testing.T) {
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
	node3, err := service.CreateNode(topoId, "test-node-3", topology.NodeTypeRouter)
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
		node3.ID,
	)
	if err != nil {
		t.Fatal(err)
	}

	testTopology_refetch1, err := service.GetTopology(testTopology.ID)
	if err != nil {
		t.Fatal(err)
	}

	linkList1, err := service.ListLinks(testTopology_refetch1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(linkList1) != 2 {
		t.Fatal(err)
	}

	service.DeleteNode(topoId, node2.ID)

	testTopology_refetch2, err := service.GetTopology(testTopology.ID)
	if err != nil {
		t.Fatal(err)
	}

	linkList2, err := service.ListLinks(testTopology_refetch2.ID)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(len(linkList2))

	if len(linkList2) != 0 {
		t.Fatal(err)
	}

}
