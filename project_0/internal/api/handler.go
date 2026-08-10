package api

import "netsim_0/internal/store"

type Handler struct {
	Store *store.MemoryStore
}

func NewHandler(s *store.MemoryStore) *Handler {
	return &Handler{
		Store: s,
	}
}
