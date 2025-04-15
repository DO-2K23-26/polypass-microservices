package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DO-2K23-26/polypass-microservices/statistics/application/services"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *services.MetricsService
}

func NewHandler(service *services.MetricsService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// User metrics endpoints
	router.HandleFunc("/api/v1/metrics/users/{userID}", h.getUserMetrics).Methods("GET")
	router.HandleFunc("/api/v1/metrics/users/{userID}/passwords/analysis", h.getPasswordStrength).Methods("GET")
	router.HandleFunc("/api/v1/metrics/users/{userID}/passwords/reused", h.getReusedPasswords).Methods("GET")
	router.HandleFunc("/api/v1/metrics/users/{userID}/passwords/breached", h.getBreachedCredentials).Methods("GET")
	router.HandleFunc("/api/v1/metrics/users/{userID}/passwords/old", h.getOldPasswords).Methods("GET")

	// Group metrics endpoints
	router.HandleFunc("/api/v1/metrics/groups/{groupID}", h.getGroupMetrics).Methods("GET")
	router.HandleFunc("/api/v1/metrics/trends", h.getCredentialTrends).Methods("GET")
	router.HandleFunc("/api/v1/metrics/credentials/{credentialID}/shared", h.getSharedCredentialStats).Methods("GET")
	router.HandleFunc("/api/v1/metrics/credentials/{credentialID}/access", h.getCredentialAccesses).Methods("GET")
}

func (h *Handler) getUserMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	metrics, err := h.service.GetUserMetrics(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(metrics)
}

func (h *Handler) getPasswordStrength(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	strengths, err := h.service.GetPasswordStrengths(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(strengths)
}

func (h *Handler) getReusedPasswords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	reused, err := h.service.GetReusedPasswords(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reused)
}

func (h *Handler) getBreachedCredentials(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	breached, err := h.service.GetBreachedCredentials(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(breached)
}

func (h *Handler) getOldPasswords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	// Default to 90 days if not specified
	daysThreshold := 90
	if days := r.URL.Query().Get("days"); days != "" {
		if d, err := time.ParseDuration(days + "d"); err == nil {
			daysThreshold = int(d.Hours() / 24)
		}
	}

	old, err := h.service.GetOldPasswords(r.Context(), userID, daysThreshold)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(old)
}

func (h *Handler) getGroupMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]

	metrics, err := h.service.GetGroupMetrics(r.Context(), groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(metrics)
}

func (h *Handler) getCredentialTrends(w http.ResponseWriter, r *http.Request) {
	startDate, err := time.Parse(time.RFC3339, r.URL.Query().Get("start"))
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0) // Default to last month
	}

	endDate, err := time.Parse(time.RFC3339, r.URL.Query().Get("end"))
	if err != nil {
		endDate = time.Now()
	}

	trends, err := h.service.GetCredentialTrends(r.Context(), startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(trends)
}

func (h *Handler) getSharedCredentialStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	credentialID := vars["credentialID"]

	stats, err := h.service.GetSharedCredentialStats(r.Context(), credentialID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) getCredentialAccesses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	credentialID := vars["credentialID"]

	startDate, err := time.Parse(time.RFC3339, r.URL.Query().Get("start"))
	if err != nil {
		startDate = time.Now().AddDate(0, -1, 0) // Default to last month
	}

	endDate, err := time.Parse(time.RFC3339, r.URL.Query().Get("end"))
	if err != nil {
		endDate = time.Now()
	}

	accesses, err := h.service.GetCredentialAccesses(r.Context(), credentialID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accesses)
}
