package api

import "netsim_0/internal/store"

type Handler struct {
	Store *store.MemoryStore
}

// Note: Memory store will eventually become an interface
// as we eventually will work with multiple different stores

func NewHandler(s *store.MemoryStore) *Handler {
	return &Handler{
		Store: s,
	}
}
