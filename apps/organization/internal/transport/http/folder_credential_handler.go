package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/DO-2K23-26/polypass-microservices/organization/internal/app"
	"github.com/gorilla/mux"
)

// FolderCredentialHandler handles folder credential endpoints.
type FolderCredentialHandler struct {
	service *app.FolderCredentialService
}

// NewFolderCredentialHandler creates a new FolderCredentialHandler.
func NewFolderCredentialHandler(svc *app.FolderCredentialService) *FolderCredentialHandler {
	return &FolderCredentialHandler{service: svc}
}

// ListCredentials handles GET /folders/{folderId}/credentials/{type}
func (h *FolderCredentialHandler) ListCredentials(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderID := vars["folderId"]
	credType := vars["type"]

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "invalid page", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	res, err := h.service.List(folderID, credType, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// CreateCredential handles POST /folders/{folderId}/credentials/{type}
func (h *FolderCredentialHandler) CreateCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderID := vars["folderId"]
	credType := vars["type"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cred, err := h.service.Create(folderID, credType, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cred)
}

// UpdateCredential handles PUT /folders/{folderId}/credentials/{type}/{credentialId}
func (h *FolderCredentialHandler) UpdateCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderID := vars["folderId"]
	credType := vars["type"]
	credentialID := vars["credentialId"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cred, err := h.service.Update(folderID, credType, credentialID, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cred)
}

// DeleteCredentials handles DELETE /folders/{folderId}/credentials/{type}
func (h *FolderCredentialHandler) DeleteCredentials(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderID := vars["folderId"]
	credType := vars["type"]
	ids := r.URL.Query()["id"]
	if len(ids) == 0 {
		http.Error(w, "missing id parameters", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(folderID, credType, ids); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
