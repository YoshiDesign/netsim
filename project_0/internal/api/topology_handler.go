package api

import (
	"encoding/json"
	"log"
	"net/http"
	"netsim_0/internal/topology"
)

// Rudimentary in-memory state per the simplicity of project-0
var Topologies = make(map[string]topology.Topology)
var nextTopologyID = 1

/*
*
 */
func (h *Handler) GetTopology(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("topoId")
	if id == "" {
		// 400 - This check is optional. Just being explicit
		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)

		return
	}

	// RLOCK
	topo, exists := h.Store.GetTopology(id)
	// RLOCK
	if !exists {
		// 404
		http.Error(
			w,
			"not found",
			http.StatusNotFound,
		)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(topo); err != nil {
		log.Printf("failed to send topology information: %v", err)
	}

}

/*
* Display all topologies
 */
func (h *Handler) GetTopologies(w http.ResponseWriter, r *http.Request) {

	// Read-Lock
	result := h.Store.ListTopologies()
	// Read-Lock

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

/*
*
 */
func (h *Handler) CreateTopology(w http.ResponseWriter, r *http.Request) {
	var req topology.CreateTopologyRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// 400
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	if req.Name == "" {
		// 400
		http.Error(
			w,
			"name is required",
			http.StatusBadRequest,
		)
		return
	}

	// Critical section
	topology := h.Store.CreateTopology(req.Name)
	// Critical section

	/*
	* HTTP/1.1 201 Created
	* Content-Type: application/json
	 */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// serialize
	if err := json.NewEncoder(w).Encode(topology); err != nil {
		log.Printf("failed to encode topology: %v", err)
	}

	// Done
}

/*
*
 */
func (h *Handler) DeleteTopology(w http.ResponseWriter, r *http.Request) {
	topoId := r.PathValue("topoId")

	if topoId == "" {
		// 400
		http.Error(
			w,
			"error: missing id",
			http.StatusBadRequest,
		)
		return
	}

	// Critical Section
	h.Store.DeleteTopology(topoId)
	// Critical Section

	// Success - We'll return a 204 No Content response
	w.WriteHeader(http.StatusNoContent)

	// ^That's all that's necessary. The header will be written
	// as part of the inevitable HTTP response.
	// We do not need to encode a reply manually
	// unless there's a specific need.
}
