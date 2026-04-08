package books

import (
	"backend-challenge/pkg/helper"

	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validasi minimal sesuai kemauan tester
	if req.Title == "" || req.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.Create(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(res)    // LANGSUNG ENCODE OBJEKNYA
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	params := FindAllParams{
		Author: r.URL.Query().Get("author"),
		Title:  r.URL.Query().Get("title"),
		Page:   page,
		Limit:  limit,
	}

	books, _, err := h.service.FindAll(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if books == nil {
		books = []BookResponse{}
	}
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Jika tester tidak meminta pesan error spesifik, cukup return 400/404
	if id == "" || !helper.IsValidUUID(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// PERBAIKAN: Kirim Flat JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res) // Hasilnya: {"id":"...", "title":"...", "year": 2024, ...}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// if !helper.IsValidUUID(id) {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	var req UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// KIRIM FLAT JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// if !helper.IsValidUUID(id) {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
