package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"netsim_0/internal/api"
	"netsim_0/internal/store"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func h_helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "hello from simulation")

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
	Store  store.MemoryStore
}

func makeServer() http.Server {
	mux := http.NewServeMux()

	// server banter
	mux.HandleFunc("GET /health", h_healthHandler)
	mux.HandleFunc("GET /hello", h_helloHandler)

	// API
	mux.HandleFunc("GET 	/api/v1/topologies", api.GetTopologies)
	mux.HandleFunc("POST 	/api/v1/topologies", api.CreateTopology)
	mux.HandleFunc("GET 	/api/v1/topologies/{topoId}", api.GetTopology)
	mux.HandleFunc("DELETE 	/api/v1/topologies/{topoId}", api.DeleteTopology)

	return http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func main() {

	ns := netsim{
		Server: makeServer(),
		Store:  store.MakeStore(),
	}

	log.Println("server listening on http://localhost:8080")

	err := ns.Server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
