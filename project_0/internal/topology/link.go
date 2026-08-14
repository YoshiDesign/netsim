package topology

type Link struct {
	ID    string `json:"id"`
	NodeA string `json:"node_a"`
	NodeB string `json:"node_b"`
}

// Normalize links
func canonicalEndpoints(nodeA, nodeB string) (string, string) {
	if nodeA > nodeB {
		return nodeB, nodeA
	}

	return nodeA, nodeB
}

// TODO - make a simple validation abstraction.
