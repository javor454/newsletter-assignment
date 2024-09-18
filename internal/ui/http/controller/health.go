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
	httpServer *http_server.Server,
	hm *healthcheck.HealthMonitor,
) *HealthController {
	controller := &HealthController{lg: lg, healthMonitor: hm}
	httpServer.GetEngine().GET("api/health/liveness", controller.Liveness)
	httpServer.GetEngine().GET("api/health/readiness", controller.Readiness)

	return controller
}

// Liveness
//	@Router		/api/health/liveness [get]
//	@Summary	Liveness - determines if app is running
//	@Tags		health
//	@Success	200
func (c *HealthController) Liveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}

// Readiness
//	@Router		/api/health/readiness [get]
//	@Summary	Readiness - determines if app is ready to receive traffic
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	response.HealthStatus
func (c *HealthController) Readiness(ctx *gin.Context) {
	const (
		healthy   = "healthy"
		unhealthy = "unhealthy"
	)

	statuses := c.healthMonitor.GetStatus()

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
