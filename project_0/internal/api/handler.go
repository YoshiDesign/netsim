package api

import (
	"errors"
	"net/http"
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

// Centralized error resolution. This also makes it impossible for
// an API call to leak internal error messages to the frontend.
func ResolveError(err error) (int, string) {
	switch {
	// Topologies
	case errors.Is(err, topology.ErrTopologyNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, topology.ErrEmptyNodeName):
		return http.StatusBadRequest, err.Error()

	// Nodes
	case errors.Is(err, topology.ErrInvalidNodeType):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, topology.ErrDuplicateNode):
		return http.StatusConflict, err.Error()

	// Links
	case errors.Is(err, topology.ErrDuplicateLink):
		return http.StatusConflict, err.Error()

	case errors.Is(err, topology.ErrLinkNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, topology.ErrEmptyLinkEndpoint):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, topology.ErrSelfLink):
		return http.StatusConflict, err.Error()

	default:
		return http.StatusInternalServerError, "internal server error"
	}

}
