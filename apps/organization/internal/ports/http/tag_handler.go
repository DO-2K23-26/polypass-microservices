package http

import (
    "encoding/json"
    "net/http"

    "github.com/DO-2K23-26/polypass-microservices/organization/internal/app"
    "github.com/DO-2K23-26/polypass-microservices/organization/internal/domain"
)

type TagHandler struct {
    service *app.TagService
}

func NewTagHandler(service *app.TagService) *TagHandler {
    return &TagHandler{service: service}
}

func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
    var tag domain.Tag
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
