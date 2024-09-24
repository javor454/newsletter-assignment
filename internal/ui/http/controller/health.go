package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/healthcheck"
	"github.com/javor454/newsletter-assignment/app/http_server"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/internal/ui/http/response"
)

type HealthController struct {
	lg            logger.Logger
	healthMonitor *healthcheck.HealthMonitor
}

func NewHealthController(
	lg logger.Logger,
	hm *healthcheck.HealthMonitor,
) *HealthController {
	controller := &HealthController{lg: lg, healthMonitor: hm}

	return controller
}

func (h *HealthController) RegisterHealhController(httpServer *http_server.Server) {
	httpServer.GetEngine().GET("api/health/liveness", h.Liveness)
	httpServer.GetEngine().GET("api/health/readiness", h.Readiness)
}

// Liveness
//
//	@Summary	Determines if app is running
//	@Router		/api/health/liveness [get]
//	@Tags		health
//
//	@Success	200
func (h *HealthController) Liveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}

// Readiness
//
//	@Summary	Determines if app is ready to receive traffic
//	@Router		/api/health/readiness [get]
//	@Tags		health
//	@Produce	json
//
//	@Success	200	{object}	response.HealthStatus
func (h *HealthController) Readiness(ctx *gin.Context) {
	const (
		healthy   = "healthy"
		unhealthy = "unhealthy"
	)

	statuses := h.healthMonitor.GetStatus()

	code := http.StatusOK
	overallStatus := healthy
	indicators := make([]response.Indicator, 0, len(statuses))
	for k, s := range statuses {
		status := healthy
		if !s {
			code = http.StatusServiceUnavailable
			status = healthy
			overallStatus = unhealthy
		}

		indicators = append(indicators, response.Indicator{
			Name:   k,
			Status: status,
		})
	}

	ctx.JSON(code, response.HealthStatus{
		Status:     overallStatus,
		Indicators: indicators,
	})
}
