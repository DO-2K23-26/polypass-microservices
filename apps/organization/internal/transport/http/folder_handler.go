package http

import (
	"encoding/json"
	"net/http"

	"github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"github.com/DO-2K23-26/polypass-microservices/organization/internal/app"
)

type FolderHandler struct {
    service *app.FolderService
}

func NewFolderHandler(service *app.FolderService) *FolderHandler {
    return &FolderHandler{service: service}
}

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
    var folder organization.Folder
    if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateFolder(folder); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}
