package api

import (
	"netsim_0/internal/topology"
)

type Handler struct {
	topologies *topology.TopologyService
}

// Note: Memory store will eventually become an interface
// as we eventually will work with multiple different stores

func NewHandler(ts *topology.TopologyService) *Handler {
	return &Handler{
		topologies: ts,
	}
}
