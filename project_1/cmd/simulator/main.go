package main

import (
	"encoding/json"
	"log"
	"net/http"
	"netsim_1/internal/api"
	"netsim_1/internal/store"
	"netsim_1/internal/topology"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func h_healthHandler(w http.ResponseWriter, r *http.Request) {

	response := HealthResponse{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
	}

}

type netsim struct {
	Server http.Server
}

func makeServer(ts *topology.TopologyService) http.Server {
	mux := http.NewServeMux()

	// Dependency injected handler
	apiHandler := api.NewHandler(ts)

	// server banter
	mux.HandleFunc("GET /health", h_healthHandler)

	// API - Topologies
	mux.HandleFunc("GET /api/v1/topologies", apiHandler.GetTopologies)
	mux.HandleFunc("POST /api/v1/topologies", apiHandler.CreateTopology)
	mux.HandleFunc("GET /api/v1/topologies/{topoId}", apiHandler.GetTopology)
	mux.HandleFunc("DELETE /api/v1/topologies/{topoId}", apiHandler.DeleteTopology)

	// API - Nodes
	mux.HandleFunc("GET /api/v1/topologies/{topoId}/nodes", apiHandler.GetNodes)
	mux.HandleFunc("POST /api/v1/topologies/{topoId}/nodes", apiHandler.CreateNode)
	mux.HandleFunc("GET /api/v1/topologies/{topoId}/nodes/{nodeId}", apiHandler.GetNode)
	mux.HandleFunc("DELETE /api/v1/topologies/{topoId}/nodes/{nodeId}", apiHandler.DeleteNode)

	// API - Links
	mux.HandleFunc("GET /api/v1/topologies/{topoId}/links", apiHandler.GetLinks)
	mux.HandleFunc("POST /api/v1/topologies/{topoId}/links", apiHandler.CreateLink)
	mux.HandleFunc("GET /api/v1/topologies/{topoId}/links/{linkId}", apiHandler.GetLink)
	mux.HandleFunc("DELETE /api/v1/topologies/{topoId}/links/{linkId}", apiHandler.DeleteLink)

	mux.HandleFunc("POST /api/v1/topologies/{topoId}/nodes/{nodeId}/interfaces", apiHandler.CreateInterface)
	mux.HandleFunc("PUT /api/v1/topologies/{topoId}/nodes/{nodeId}/interfaces/{interfaceId}/address", apiHandler.SetInterfaceAddress)

	return http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func main() {

	// Dep's

	/**
	* Repository layer - owns state correcntess
	 */
	store := store.MakeStore()

	/*
	* Service layer - injects repository
	* Owns domain correctness
	 */
	topologyService := topology.NewTopologyService(store)

	/**
	* Constructs the API Layer
	* API Layer owns transport correcntess
	 */
	ns := netsim{
		Server: makeServer(topologyService), // inject topology service
	}

	log.Println("server listening on http://localhost:8080")

	err := ns.Server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
