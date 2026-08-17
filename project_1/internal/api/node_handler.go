package api

import (
	"encoding/json"
	"log"
	"net/http"
	"netsim_1/internal/topology"
)

type createNodeRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (a *Handler) CreateNode(w http.ResponseWriter, r *http.Request) {

	topologyId := r.PathValue("topoId")

	/*
	* Note: Node names & ID will be validated
	* in the service layer.
	 */

	var req createNodeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	node, err := a.topologies.CreateNode(
		topologyId,
		topology.NodeName(req.Name),
		topology.NodeType(req.Type),
	)

	// Error Handling - Every type that appears throughout the CreateNode process.
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

	// Success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(node); err != nil {
		// Log failures
		log.Printf("failed to encode response: %v", err)
	}
}

func (a *Handler) GetNodes(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")
	nodes, err := a.topologies.ListNodes(topologyId)

	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(nodes); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (a *Handler) GetNode(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")
	nodeId := r.PathValue("nodeId")

	node, err := a.topologies.GetNode(topologyId, topology.NodeID(nodeId))
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(node); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (a *Handler) DeleteNode(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")
	nodeId := r.PathValue("nodeId")

	err := a.topologies.DeleteNode(topologyId, topology.NodeID(nodeId))
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

	// 204
	w.WriteHeader(http.StatusNoContent)

	// Alternatively
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("deleted"))

	// Or
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// json.NewEncoder(w).Encode(map[string]int{
	// 	"deleted": 1,
	// })

	// Or define a DeleteResponse

}
