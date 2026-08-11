package topology

import (
	"errors"
	"strings"
)

/*
* The service layer enforces the domain rules and policies
* before forwarding requests to the API layer.
*
* Note the use of errors.Is() something I learned here is that
* this works even if returned values are via fmt.Errorf. Pretty cool
 */

// Domain specific (non-HTTP) errors for the service layer
var (
	ErrTopologyNotFound     = errors.New("topology not found")
	ErrTopologyNameRequired = errors.New("topology name is required")
	ErrTopologyIdNotFound   = errors.New("topology ID not found")
)

/* */
func NewTopologyService(store TopologyStore) *TopologyService {
	return &TopologyService{
		store: store,
	}
}

/* */
func (s *TopologyService) CreateTopology(name string) (Topology, error) {
	// Basic name validation
	name = strings.TrimSpace(name)
	if name == "" {
		return Topology{}, ErrTopologyNameRequired
	}

	created := s.store.CreateTopology(name)

	return created, nil
}

/* */
func (s *TopologyService) ListTopologies() []Topology {
	return s.store.ListTopologies()
}

/* */
func (s *TopologyService) GetTopology(id string) (Topology, error) {
	t, exists := s.store.GetTopology(id)
	if !exists {
		return Topology{}, ErrTopologyNotFound
	}

	return t, nil
}

/* */
func (s *TopologyService) DeleteTopology(id string) error {
	deleted := s.store.DeleteTopology(id)
	if !deleted {
		return ErrTopologyNotFound
	}

	return nil
}
