package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.yurtal.tech/company/blitz/back/internal/model"
	"gitlab.yurtal.tech/company/blitz/back/internal/service"
)

// CreateDriver godoc
// @Summary Create a new driver
// @Description Create a new driver with the provided information
// @Tags drivers
// @Accept json
// @Produce json
// @Param driver body model.CreateDriverRequest true "Driver information"
// @Success 201 {object} model.DriverResponse
// @Failure 400 {object} model.DriverErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 409 {object} model.DriverErrorResponse
// @Failure 500 {object} model.DriverErrorResponse
// @Security BearerAuth
// @Router /api/v1/drivers [post]
func (h *Handler) CreateDriver(c echo.Context) error {
	var req model.CreateDriverRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.DriverErrorResponse{
			Error: model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "invalid request body",
			},
		})
	}

	driver, err := h.service.Driver().CreateDriver(c.Request().Context(), req)
	if err != nil {
		return h.handleDriverError(c, err)
	}

	return c.JSON(http.StatusCreated, driver)
}

// GetDriver godoc
// @Summary Get driver by ID
// @Description Get detailed information about a specific driver
// @Tags drivers
// @Produce json
// @Param id path string true "Driver ID (UUID)"
// @Success 200 {object} model.DriverResponse
// @Failure 400 {object} model.DriverErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.DriverErrorResponse
// @Failure 500 {object} model.DriverErrorResponse
// @Security BearerAuth
// @Router /api/v1/drivers/{id} [get]
func (h *Handler) GetDriver(c echo.Context) error {
	id := c.Param("id")

	driver, err := h.service.Driver().GetDriver(c.Request().Context(), id)
	if err != nil {
		return h.handleDriverError(c, err)
	}

	return c.JSON(http.StatusOK, driver)
}

// ListDrivers godoc
// @Summary List all drivers
// @Description Get a paginated list of drivers with optional filters
// @Tags drivers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query model.DriverStatus false "Filter by status"
// @Param search query string false "Search by name or phone"
// @Success 200 {object} model.ListDriversResponse
// @Failure 400 {object} model.DriverErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.DriverErrorResponse
// @Security BearerAuth
// @Router /api/v1/drivers [get]
func (h *Handler) ListDrivers(c echo.Context) error {
	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	status := c.QueryParam("status")
	search := c.QueryParam("search")

	response, err := h.service.Driver().ListDrivers(c.Request().Context(), page, limit, status, search)
	if err != nil {
		return h.handleDriverError(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateDriver godoc
// @Summary Update driver information
// @Description Partially update driver information
// @Tags drivers
// @Accept json
// @Produce json
// @Param id path string true "Driver ID (UUID)"
// @Param driver body model.UpdateDriverRequest true "Fields to update"
// @Success 200 {object} model.DriverResponse
// @Failure 400 {object} model.DriverErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.DriverErrorResponse
// @Failure 409 {object} model.DriverErrorResponse
// @Failure 500 {object} model.DriverErrorResponse
// @Security BearerAuth
// @Router /api/v1/drivers/{id} [patch]
func (h *Handler) UpdateDriver(c echo.Context) error {
	id := c.Param("id")

	var req model.UpdateDriverRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.DriverErrorResponse{
			Error: model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "invalid request body",
			},
		})
	}

	driver, err := h.service.Driver().UpdateDriver(c.Request().Context(), id, req)
	if err != nil {
		return h.handleDriverError(c, err)
	}

	return c.JSON(http.StatusOK, driver)
}

// UpdateDriverStatus godoc
// @Summary Update driver status
// @Description Change driver status (active, inactive, blocked)
// @Tags drivers
// @Accept json
// @Produce json
// @Param id path string true "Driver ID (UUID)"
// @Param status body model.UpdateStatusRequest true "New status"
// @Success 200 {object} model.DriverResponse
// @Failure 400 {object} model.DriverErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.DriverErrorResponse
// @Failure 500 {object} model.DriverErrorResponse
// @Security BearerAuth
// @Router /api/v1/drivers/{id}/status [patch]
func (h *Handler) UpdateDriverStatus(c echo.Context) error {
	id := c.Param("id")

	var req model.UpdateStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.DriverErrorResponse{
			Error: model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "invalid request body",
			},
		})
	}

	driver, err := h.service.Driver().UpdateDriverStatus(c.Request().Context(), id, req)
	if err != nil {
		return h.handleDriverError(c, err)
	}

	return c.JSON(http.StatusOK, driver)
}

// DeleteDriver godoc
// @Summary Delete driver
// @Description Soft delete a driver (marks as deleted but retains data)
// @Tags drivers
// @Param id path string true "Driver ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} model.DriverErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.DriverErrorResponse
// @Failure 500 {object} model.DriverErrorResponse
// @Security BearerAuth
// @Router /api/v1/drivers/{id} [delete]
func (h *Handler) DeleteDriver(c echo.Context) error {
	id := c.Param("id")

	err := h.service.Driver().DeleteDriver(c.Request().Context(), id)
	if err != nil {
		return h.handleDriverError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) handleDriverError(c echo.Context, err error) error {
	if appErr, ok := err.(service.AppError); ok {
		status := http.StatusBadRequest

		switch appErr.Code {
		case "DRIVER_NOT_FOUND":
			status = http.StatusNotFound
		case "PHONE_CONFLICT", "LICENSE_CONFLICT", "CAR_PLATE_CONFLICT":
			status = http.StatusConflict
		case "VALIDATION_ERROR", "INVALID_STATUS", "INVALID_ID":
			status = http.StatusBadRequest
		}

		return c.JSON(status, model.DriverErrorResponse{
			Error: model.ErrorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
	}

	h.logger.Error("internal server error", "error", err)
	return c.JSON(http.StatusInternalServerError, model.DriverErrorResponse{
		Error: model.ErrorDetail{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		},
	})
}

// GetActiveDriversStats godoc
// @Summary Get driver status statistics
// @Description Get count of drivers by status
// @Tags drivers
// @Produce json
// @Param status query model.DriverStatus false "Driver status" default(active)
// @Success 200 {object} model.ActiveDriversStatsResponse
// @Failure 400 {object} model.DriverErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.DriverErrorResponse
// @Security BearerAuth
// @Router /api/v1/drivers/stats/active [get]
func (h *Handler) GetActiveDriversStats(c echo.Context) error {
	status := c.QueryParam("status")

	stats, err := h.service.Driver().GetActiveDriversStats(c.Request().Context(), status)
	if err != nil {
		return h.handleDriverError(c, err)
	}

	return c.JSON(http.StatusOK, stats)
}
