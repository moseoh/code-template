package common_handler

import (
	"encoding/json"
	"net/http"

	"pkg/models"
)

type HealthHandler struct {
	serviceName string
	version     string
}

func NewHealthHandler(serviceName, version string) *HealthHandler {
	return &HealthHandler{
		serviceName: serviceName,
		version:     version,
	}
}

// Path 헬스체크 경로 반환
func (h *HealthHandler) Path() string {
	return "/health"
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := models.NewHealthResponse(h.serviceName, h.version)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
