package api

import (
	"encoding/json"
	"log"
	"net/http"
	"netsim_1/internal/topology"
)

type CreateIfaceRequest struct {
	Name topology.InterfaceName `json:"name"`
}

type UpdateInterfaceRequest struct {
	Address string `json:"address"`
}

func (a *Handler) CreateInterface(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")
	nodeId := topology.NodeID(r.PathValue("nodeId"))

	var req CreateIfaceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	iface, err := a.topologies.CreateInterface(topologyId, nodeId, req.Name)
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(iface); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}

func (a *Handler) SetInterfaceAddress(w http.ResponseWriter, r *http.Request) {
	topologyId := r.PathValue("topoId")
	nodeId := topology.NodeID(r.PathValue("nodeId"))
	ifaceId := topology.InterfaceID(r.PathValue("interfaceId"))

	var req UpdateInterfaceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	iface, err := a.topologies.SetInterfaceAddress(
		topologyId,
		nodeId,
		ifaceId,
		req.Address,
	)
	if err != nil {
		status, message := ResolveError(err)
		http.Error(w, message, status)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(iface); err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}
