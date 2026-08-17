package topology

import "strings"

type LinkName string
type LinkID string

type Link struct {
	ID        LinkID   `json:"id"`
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

func (s LinkName) TrimSpace() LinkName {
	return LinkName(strings.TrimSpace(string(s)))
}

func (s LinkID) TrimSpace() LinkID {
	return LinkID(strings.TrimSpace(string(s)))
}

// TODO - make a simple validation abstraction.
