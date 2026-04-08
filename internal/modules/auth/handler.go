package auth

import (
	"backend-challenge/pkg/response"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GenerateToken()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed generate token")
		return
	}
	response.Success(w, http.StatusOK, "token generated", res)
}
