package api

import (
	"net/http"
	"homelab.com/homelab-server/homeLab-server/internal/repositories"
	"encoding/json"
)

type MetricsHandler struct {
	Repo repositories.PrometheusRepository
}

func NewMetricsHandler(repo repositories.PrometheusRepository) *MetricsHandler {
	return &MetricsHandler{Repo: repo}
}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query") // Get the query parameter from the request
	metrics, err := h.Repo.GetMetrics(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
} 