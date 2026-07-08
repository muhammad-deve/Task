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
	Code          int          `json:"code,omitempty"`
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
	Code int              `json:"code"`
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
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ActiveDriversStatsResponse struct {
	Code   int          `json:"code"`
	Status DriverStatus `json:"status"`
	Count  int          `json:"count"`
}
