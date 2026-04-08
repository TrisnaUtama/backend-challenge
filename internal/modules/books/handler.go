package books

import (
	"backend-challenge/pkg/helper"
	"backend-challenge/pkg/response"

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
		response.ValidationError(w, map[string]string{
			"body": "invalid JSON format",
		})
		return
	}

	errors := map[string]string{}
	if req.Title == "" {
		errors["title"] = "title is required"
	}
	if req.Author == "" {
		errors["author"] = "author is required"
	}
	if len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	res, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "book created", res)
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

	books, total, err := h.service.FindAll(r.Context(), params)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	response.Paginated(w, http.StatusOK, "success", books, response.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	})
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if !helper.IsValidUUID(id) {
		response.ValidationError(w, map[string]string{
			"id": "invalid UUID format",
		})
		return
	}

	if id == "" {
		response.ValidationError(w, map[string]string{
			"id": "id is required",
		})
		return
	}

	res, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			response.Error(w, http.StatusNotFound, "book not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "success", res)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if !helper.IsValidUUID(id) {
		response.ValidationError(w, map[string]string{
			"id": "invalid UUID format",
		})
		return
	}

	var req UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, map[string]string{
			"body": "invalid JSON",
		})
		return
	}

	errs := map[string]string{}
	if req.Title == "" {
		errs["title"] = "title is required"
	}
	if req.Author == "" {
		errs["author"] = "author is required"
	}
	if len(errs) > 0 {
		response.ValidationError(w, errs)
		return
	}

	res, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			response.Error(w, http.StatusNotFound, "book not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "book updated", res)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if !helper.IsValidUUID(id) {
        response.ValidationError(w, map[string]string{
            "id": "invalid UUID format",
        })
        return
    }

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if err == ErrBookNotFound {
			response.Error(w, http.StatusNotFound, "book not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "book deleted", nil)
}
