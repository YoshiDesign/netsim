package topology

type NodeType string

const (
	NodeTypeRouter NodeType = "router"
	NodeTypeSwitch NodeType = "switch"
	NodeTypeHost   NodeType = "host"
)

type Node struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
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
