package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"netsim_0/internal/topology"
)

// Rudimentary in-memory state per the simplicity of project-0
var Topologies = make(map[string]topology.Topology)
var nextTopologyID = 1

func GetTopology(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("topoId")
	if id == "" {
		// 400 - could alternatively drop the request silently. Design/security detail
		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)

		return
	}

	// Get from global memory
	topo, exists := Topologies[id]
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

// Display all topologies
func GetTopologies(w http.ResponseWriter, r *http.Request) {

	result := make([]topology.Topology, 0, len(Topologies))

	for _, topo := range Topologies {
		result = append(result, topo)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func CreateTopology(w http.ResponseWriter, r *http.Request) {
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

	id := fmt.Sprintf("topology-%d", nextTopologyID)
	nextTopologyID++

	// DTO
	topology := topology.Topology{
		ID:   id,
		Name: req.Name,
	}

	// "Register" the new topology
	Topologies[id] = topology

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

func DeleteTopology(w http.ResponseWriter, r *http.Request) {
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

	_, exists := Topologies[topoId]
	if !exists {
		// 404
		http.Error(
			w,
			"error: missing id",
			http.StatusNotFound,
		)
		return
	}

	delete(Topologies, topoId)

	// Success - We'll return a 204 No Content response
	w.WriteHeader(http.StatusNoContent)

	// That's all that's necessary. The header will be written
	// as part of the inevitable HTTP response.
	// We do not need to encode a reply manually
	// unless there's a specific need.
}
