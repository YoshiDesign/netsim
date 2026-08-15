package topology

import "strings"

type NodeType string
type NodeName string
type NodeID string

const INVALID_NODE = "00000000000X" // or some unique hash preferably, once things progress

func (s NodeName) TrimSpace() NodeName {
	return NodeName(strings.TrimSpace(string(s)))
}

func (s NodeID) TrimSpace() NodeID {
	return NodeID(strings.TrimSpace(string(s)))
}

// TODO func (s NodeName) Valid() NodeName etc.

const (
	NodeTypeRouter NodeType = "router"
	NodeTypeSwitch NodeType = "switch"
	NodeTypeHost   NodeType = "host"
)

type Node struct {
	ID   NodeID   `json:"id"`
	Name NodeName `json:"name"`
	Type NodeType `json:"type"`
}

func (t NodeType) Valid() bool {
	switch t {
	case NodeTypeRouter, NodeTypeSwitch, NodeTypeHost:
		return true
	default:
		return false
	}
}
