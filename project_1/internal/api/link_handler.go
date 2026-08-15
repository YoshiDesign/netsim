package api

import (
	"encoding/json"
	"log"
	"net/http"
	"netsim_1/internal/topology"
)

type CreateLinkRequest struct {
	NodeA topology.NodeID `json:"node_a"`
	NodeB topology.NodeID `json:"node_b"`
}

/*
* Create a link between nodes
* See: CreateLinkRequest ^
 */
func (a *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")

	var req CreateLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	link, err := a.topologies.CreateLink(topologyId, topology.NodeID(req.NodeA), topology.NodeID(req.NodeB))
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(link); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (a *Handler) GetLinks(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")

	links, err := a.topologies.ListLinks(topologyId)
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(links); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}

func (a *Handler) GetLink(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")
	linkId := r.PathValue("linkId")

	link, err := a.topologies.GetLink(topologyId, linkId)
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(link); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (a *Handler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")
	linkId := r.PathValue("linkId")

	err := a.topologies.DeleteLink(topologyId, linkId)
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

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
