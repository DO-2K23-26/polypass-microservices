package http

import (
	"encoding/json"
	"net/http"

	"github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/app"
	"github.com/gorilla/mux"
)

type TagHandler struct {
	service *app.TagService
}

func NewTagHandler(service *app.TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var tag organization.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTag(tag); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var tag organization.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tag.Id = id
	if err := h.service.UpdateTag(tag); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.service.DeleteTag(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.service.ListTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func (h *TagHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tag, err := h.service.GetTag(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tag)
}
