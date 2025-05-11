package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/polypass/polypass-microservices/statistics/application/services"
	"github.com/polypass/polypass-microservices/statistics/domain/models"
)

// RestHandler handles HTTP requests for the statistics service
type RestHandler struct {
	metricService *services.MetricService
	eventService  *services.EventService
}

// NewRestHandler creates a new instance of RestHandler
func NewRestHandler(metricService *services.MetricService, eventService *services.EventService) *RestHandler {
	return &RestHandler{
		metricService: metricService,
		eventService:  eventService,
	}
}

// RegisterRoutes registers all the routes for the statistics service
func (h *RestHandler) RegisterRoutes(router *mux.Router) {
	// Metrics endpoints
	router.HandleFunc("/metrics", h.GetLatestMetrics).Methods("GET")
	router.HandleFunc("/metrics/{id}", h.GetMetricByID).Methods("GET")
	router.HandleFunc("/metrics/name/{name}", h.GetMetricsByName).Methods("GET")
	router.HandleFunc("/metrics/category/{category}", h.GetMetricsByCategory).Methods("GET")
	router.HandleFunc("/metrics/calculate", h.CalculateMetrics).Methods("POST")
}

// GetLatestMetrics handles requests to get the latest metrics
func (h *RestHandler) GetLatestMetrics(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // Default limit

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	metrics, err := h.metricService.GetLatestMetrics(r.Context(), limit)
	if err != nil {
		http.Error(w, "Failed to get metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, metrics)
}

// GetMetricByID handles requests to get a metric by ID
func (h *RestHandler) GetMetricByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	metric, err := h.metricService.GetMetricByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get metric: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, metric)
}

// GetMetricsByName handles requests to get metrics by name
func (h *RestHandler) GetMetricsByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	metrics, err := h.metricService.GetMetricsByName(r.Context(), name)
	if err != nil {
		http.Error(w, "Failed to get metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, metrics)
}

// GetMetricsByCategory handles requests to get metrics by category
func (h *RestHandler) GetMetricsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryStr := vars["category"]

	category := models.MetricCategory(categoryStr)

	metrics, err := h.metricService.GetMetricsByCategory(r.Context(), category)
	if err != nil {
		http.Error(w, "Failed to get metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, metrics)
}

// CalculateMetrics handles requests to calculate metrics
func (h *RestHandler) CalculateMetrics(w http.ResponseWriter, r *http.Request) {
	var request struct {
		TimeRange string `json:"timeRange"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	duration, err := time.ParseDuration(request.TimeRange)
	if err != nil {
		http.Error(w, "Invalid time range format", http.StatusBadRequest)
		return
	}

	metrics, err := h.metricService.CalculateMetrics(r.Context(), duration)
	if err != nil {
		http.Error(w, "Failed to calculate metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, metrics)
}

// respondWithJSON is a helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
