package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

// HealthCheck godoc
// @Summary Health check
// @Description Check service health and database connectivity
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /healthz [get]
func (h *Handler) HealthCheck(c echo.Context) error {
	ctx := c.Request().Context()

	dbStatus := "healthy"
	if err := h.pool.Ping(ctx); err != nil {
		dbStatus = "unhealthy"
	}

	response := HealthResponse{
		Status: "ok",
		DB:     dbStatus,
	}

	statusCode := http.StatusOK
	if dbStatus != "healthy" {
		statusCode = http.StatusServiceUnavailable
		response.Status = "degraded"
	}

	return c.JSON(statusCode, response)
}
