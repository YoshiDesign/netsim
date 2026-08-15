package topology

type Link struct {
	ID        string   `json:"id"`
	NodeAName NodeName `json:"node_a"`
	NodeBName NodeName `json:"node_b"`
	NodeAID   NodeID   `json:"node_a_id"`
	NodeBID   NodeID   `json:"node_b_id"`
}

/*
* Architecture invariant:
* Every link endpoint must always refer to a node that currently exists in the same topology.
* This is why deleting a node cascades into the deletion of every link it's involved in.
 */

// Normalize links
func canonicalEndpointsByID(nodeA, nodeB NodeID) (NodeID, NodeID) {
	if nodeA > nodeB {
		return nodeB, nodeA
	}

	return nodeA, nodeB
}

// TODO - make a simple validation abstraction.
