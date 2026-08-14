package topology

type Link struct {
	ID    string `json:"id"`
	NodeA string `json:"node_a"`
	NodeB string `json:"node_b"`
}

/*
* Architecture invariant:
* Every link endpoint must always refer to a node that currently exists in the same topology.
* This is why deleting a node cascades into the deletion of every link it's involved in.
 */

// Normalize links
func canonicalEndpoints(nodeA, nodeB string) (string, string) {
	if nodeA > nodeB {
		return nodeB, nodeA
	}

	return nodeA, nodeB
}

// TODO - make a simple validation abstraction.
