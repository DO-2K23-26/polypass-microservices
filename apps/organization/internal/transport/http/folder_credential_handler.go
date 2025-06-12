package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	organization "github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
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

	if organization.CredentialType(credType) != organization.CredentialTypePassword &&
		organization.CredentialType(credType) != organization.CredentialTypeCard &&
		organization.CredentialType(credType) != organization.CredentialTypeSSHKey {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errBody := map[string]interface{}{"error": "invalid credential type", "types": []string{"password", "card", "sshkey"}}
		json.NewEncoder(w).Encode(errBody)
		return
	}

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorBody := map[string]string{"error": "invalid page number"}
		json.NewEncoder(w).Encode(errorBody)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorBody := map[string]string{"error": "invalid limit number"}
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	req := organization.GetCredentialRequest{
		Page:  page,
		Limit: limit,
	}

	credTypeStr := credType
	res, err := h.service.List(folderID, &credTypeStr, &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errBody := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errBody)
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

	if organization.CredentialType(credType) != organization.CredentialTypePassword &&
		organization.CredentialType(credType) != organization.CredentialTypeCard &&
		organization.CredentialType(credType) != organization.CredentialTypeSSHKey {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errBody := map[string]interface{}{"error": "invalid credential type", "types": []string{"password", "card", "sshkey"}}
		json.NewEncoder(w).Encode(errBody)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errBody := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errBody)
		return
	}

	cred, err := h.service.Create(folderID, credType, body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errBody := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errBody)
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

	if organization.CredentialType(credType) != organization.CredentialTypePassword &&
		organization.CredentialType(credType) != organization.CredentialTypeCard &&
		organization.CredentialType(credType) != organization.CredentialTypeSSHKey {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errBody := map[string]interface{}{"error": "invalid credential type", "types": []string{"password", "card", "sshkey"}}
		json.NewEncoder(w).Encode(errBody)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errBody := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errBody)
		return
	}

	cred, err := h.service.Update(folderID, credType, credentialID, body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errBody := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errBody)
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

	if organization.CredentialType(credType) != organization.CredentialTypePassword &&
		organization.CredentialType(credType) != organization.CredentialTypeCard &&
		organization.CredentialType(credType) != organization.CredentialTypeSSHKey {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errBody := map[string]interface{}{"error": "invalid credential type", "types": []string{"password", "card", "sshkey"}}
		json.NewEncoder(w).Encode(errBody)
		return
	}

	ids := r.URL.Query()["id"]
	if len(ids) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errBody := map[string]string{"error": "missing id parameters"}
		json.NewEncoder(w).Encode(errBody)
		return
	}

	if err := h.service.Delete(folderID, credType, ids); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errBody := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errBody)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *FolderCredentialHandler) ListUserCredentials(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorBody := map[string]string{"error": "missing user_id parameter"}
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	res, err := h.service.ListUserCredentials(userId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errBody := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errBody)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
