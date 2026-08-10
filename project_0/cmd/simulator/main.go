package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"netsim_0/internal/api"
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

func main() {

	mux := http.NewServeMux()

	// server banter
	mux.HandleFunc("GET /health", h_healthHandler)
	mux.HandleFunc("GET /hello", h_helloHandler)

	// API
	mux.HandleFunc("GET 	/api/v1/topologies", api.GetTopologies)
	mux.HandleFunc("POST 	/api/v1/topologies", api.CreateTopology)
	mux.HandleFunc("GET 	/api/v1/topologies/{topoId}", api.GetTopology)
	mux.HandleFunc("DELETE 	/api/v1/topologies/{topoId}", api.DeleteTopology)

	log.Println("server listening on http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}
