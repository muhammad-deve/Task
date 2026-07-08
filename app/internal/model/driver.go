package model

import "time"

type DriverStatus string

const (
	StatusActive   DriverStatus = "active"
	StatusInactive DriverStatus = "inactive"
	StatusBlocked  DriverStatus = "blocked"
)

type CreateDriverRequest struct {
	FullName      string       `json:"full_name"`
	Phone         string       `json:"phone"`
	LicenseNumber string       `json:"license_number"`
	CarModel      string       `json:"car_model"`
	CarPlate      string       `json:"car_plate"`
	Status        DriverStatus `json:"status"`
}

type UpdateDriverRequest struct {
	FullName      *string `json:"full_name"`
	Phone         *string `json:"phone"`
	LicenseNumber *string `json:"license_number"`
	CarModel      *string `json:"car_model"`
	CarPlate      *string `json:"car_plate"`
}

type UpdateStatusRequest struct {
	Status DriverStatus `json:"status"`
}

type DriverResponse struct {
	ID            string       `json:"id"`
	FullName      string       `json:"full_name"`
	Phone         string       `json:"phone"`
	LicenseNumber string       `json:"license_number"`
	CarModel      string       `json:"car_model"`
	CarPlate      string       `json:"car_plate"`
	Status        DriverStatus `json:"status"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type ListDriversResponse struct {
	Data []DriverResponse `json:"data"`
	Meta PaginationMeta   `json:"meta"`
}

type PaginationMeta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type DriverErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ActivityAction string

const (
	ActionWentOnline  ActivityAction = "went_online"
	ActionWentOffline ActivityAction = "went_offline"
)

type LogActivityRequest struct {
	Action ActivityAction `json:"action"`
	Notes  *string        `json:"notes"`
}

type ActivityLogResponse struct {
	ID        string         `json:"id"`
	DriverID  string         `json:"driver_id"`
	Action    ActivityAction `json:"action"`
	Timestamp time.Time      `json:"timestamp"`
	Notes     *string        `json:"notes,omitempty"`
}

type WorkingHoursResponse struct {
	DriverID    string  `json:"driver_id"`
	TotalHours  float64 `json:"total_hours"`
	PeriodStart string  `json:"period_start"`
	PeriodEnd   string  `json:"period_end"`
}

type ActiveDriversStatsResponse struct {
	ActiveCount int `json:"active_count"`
}
