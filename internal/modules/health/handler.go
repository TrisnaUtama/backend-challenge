package health

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	res := h.service.Ping()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Echo(w http.ResponseWriter, r *http.Request) {
	var req interface{}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil || req == nil {
		req = map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}
