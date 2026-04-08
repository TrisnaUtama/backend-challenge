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
	var req EchoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	res := h.service.Echo(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
