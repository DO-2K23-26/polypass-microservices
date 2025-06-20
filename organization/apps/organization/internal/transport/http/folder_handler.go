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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorBody := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	folder, serviceError := h.service.CreateFolder(folderRequest)

	if serviceError != nil {
		errorBody := map[string]string{"error": serviceError.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(folder)
}

func (h *FolderHandler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderId := vars["id"]

	var folder organization.UpdateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		errorBody := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	result, serviceErr := h.service.UpdateFolder(folderId, folder)

	if serviceErr != nil {
		errorBody := map[string]string{"error": serviceErr.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.service.DeleteFolder(id); err != nil {
		errorBody := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
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

	userIdStr := r.URL.Query().Get("user_id")
	var userId *string
	if userIdStr != "" {
		userId = &userIdStr
	} else {
		userId = nil
	}

	// Convert page and limit to integers
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		errorBody := map[string]string{"error": "Invalid page number"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		errorBody := map[string]string{"error": "Invalid limit number"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}
	req := organization.GetFolderRequest{
		Page:   pageInt,
		Limit:  limitInt,
		UserId: userId,
	}

	folders, err := h.service.ListFolders(req)
	if err != nil {
		errorBody := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
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
		errorBody := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folder)
}

func (h *FolderHandler) ListUsersInFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderId := vars["id"]

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
		errorBody := map[string]string{"error": "Invalid page number"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		errorBody := map[string]string{"error": "Invalid limit number"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}
	req := organization.GetUsersInFolderRequest{
		Page:  pageInt,
		Limit: limitInt,
	}

	users, err := h.service.ListUsersInFolder(folderId, req)
	if err != nil {
		errorBody := map[string]string{"error": err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	result := organization.GetUsersInFolderResponse{
		Users: users,
		Total: len(users),
		Page:  pageInt,
		Limit: limitInt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
