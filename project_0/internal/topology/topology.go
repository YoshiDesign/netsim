package topology

// Describe what the API accepts
type CreateTopologyRequest struct {
	Name string `json:"name"`
}

// Describe a topology we already have
type Topology struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
