package http

import (
	"net/http"

	"github.com/DO-2K23-26/polypass-microservices/search-service/services/health"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	HealthService *health.HealthService
}

func NewHealthController(healthService *health.HealthService) *HealthController {
	if healthService == nil {
		panic("healthService cannot be nil")
	}
	return &HealthController{
		HealthService: healthService,
	}
}

func (hc *HealthController) CheckHealth(c *gin.Context) {
	healthStatus := hc.HealthService.CheckHealth()
	c.JSON(http.StatusOK, healthStatus)
}
