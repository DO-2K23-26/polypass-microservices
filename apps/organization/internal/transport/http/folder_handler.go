package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/app"
	"github.com/gorilla/mux"
)

type FolderHandler struct {
	service *app.FolderService
}

func NewFolderHandler(service *app.FolderService) *FolderHandler {
	return &FolderHandler{service: service}
}

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	var folderRequest organization.CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&folderRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	folder, serviceError := h.service.CreateFolder(folderRequest)

	if serviceError != nil {
		http.Error(w, serviceError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folder)
}

func (h *FolderHandler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var folder organization.Folder
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	folder.Id = id
	if err := h.service.UpdateFolder(folder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.service.DeleteFolder(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FolderHandler) ListFolders(w http.ResponseWriter, r *http.Request) {
	// Pagination parameters
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1" // Default to page 1 if not provided
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10" // Default to 10 items per page if not provided
	}
	// Convert page and limit to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit number", http.StatusBadRequest)
		return
	}
	req := organization.GetFolderRequest{
		Page:  pageInt,
		Limit: limitInt,
	}

	folders, err := h.service.ListFolders(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := organization.GetFolderResponse{
		Folders: folders,
		Total:   len(folders), // Assuming total is the length of the returned folders
		Page:    pageInt,
		Limit:   limitInt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *FolderHandler) GetFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	folder, err := h.service.GetFolder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}
