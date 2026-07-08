package handler

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"gitlab.yurtal.tech/company/blitz/back/internal/config"
	mw "gitlab.yurtal.tech/company/blitz/back/internal/middleware"
	"gitlab.yurtal.tech/company/blitz/back/internal/service"
	"gitlab.yurtal.tech/company/blitz/back/pkg/logger"
)

type Handler struct {
	logger  *logger.Logger
	service service.I
	cfg     *config.Config
	pool    *pgxpool.Pool
}

func (h *Handler) Register(router *echo.Echo) {
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(router)

	router.GET("/healthz", h.HealthCheck)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Login, mw.LoginRateLimiter(), mw.ValidateLoginInput)
			auth.POST("/register", h.RegisterUser, mw.ValidateRegisterInput)
			auth.POST("/refresh", h.Refresh, mw.ValidateRefreshInput)
		}

		drivers := api.Group("/drivers")
		{
			drivers.POST("", h.CreateDriver)
			drivers.GET("", h.ListDrivers)
			drivers.GET("/stats/active", h.GetActiveDriversStats)
			drivers.GET("/:id", h.GetDriver)
			drivers.PATCH("/:id", h.UpdateDriver)
			drivers.DELETE("/:id", h.DeleteDriver)
			drivers.PATCH("/:id/status", h.UpdateDriverStatus)
			drivers.POST("/:id/activity", h.LogDriverActivity)
			drivers.GET("/:id/activity", h.GetDriverActivityLog)
			drivers.GET("/:id/working-hours", h.GetDriverWorkingHours)
		}
	}

}

func New(logger *logger.Logger, cfg *config.Config, service service.I, pool *pgxpool.Pool) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
		cfg:     cfg,
		pool:    pool,
	}
}
