package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gitlab.yurtal.tech/company/blitz/back/internal/config"
	"gitlab.yurtal.tech/company/blitz/back/internal/model"
	"gitlab.yurtal.tech/company/blitz/back/internal/repository"
	"gitlab.yurtal.tech/company/blitz/back/internal/repository/pg"
)

type DriverService struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewDriverService(cfg *config.Config, repo *repository.Repository) *DriverService {
	return &DriverService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *DriverService) CreateDriver(ctx context.Context, req model.CreateDriverRequest) (model.DriverResponse, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return model.DriverResponse{}, err
	}

	if err := s.checkUniqueFields(ctx, "", req.Phone, req.LicenseNumber, req.CarPlate); err != nil {
		return model.DriverResponse{}, err
	}

	status := string(model.StatusActive)
	if req.Status != "" {
		status = string(req.Status)
	}

	driver, err := s.repo.PgRepo.Repo.CreateDriver(ctx, pg.CreateDriverParams{
		FullName:      req.FullName,
		Phone:         req.Phone,
		LicenseNumber: req.LicenseNumber,
		CarModel:      req.CarModel,
		CarPlate:      req.CarPlate,
		Status:        status,
	})
	if err != nil {
		return model.DriverResponse{}, fmt.Errorf("failed to create driver: %w", err)
	}

	return s.toDriverResponse(driver), nil
}

func (s *DriverService) GetDriver(ctx context.Context, id string) (model.DriverResponse, error) {
	driverID, err := uuid.Parse(id)
	if err != nil {
		return model.DriverResponse{}, NewAppError("INVALID_ID", "invalid driver id format")
	}

	driver, err := s.repo.PgRepo.Repo.GetDriverByID(ctx, driverID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DriverResponse{}, NewAppError("DRIVER_NOT_FOUND", fmt.Sprintf("driver with id %s not found", id))
		}
		return model.DriverResponse{}, fmt.Errorf("failed to get driver: %w", err)
	}

	return s.toDriverResponse(driver), nil
}

func (s *DriverService) ListDrivers(ctx context.Context, page, limit int, status, search string) (model.ListDriversResponse, error) {
	status = strings.TrimSpace(status)
	search = strings.TrimSpace(search)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	var statusFilter string
	if status != "" {
		if !isValidStatus(status) {
			return model.ListDriversResponse{}, NewAppError("INVALID_STATUS", "status must be one of: active, inactive, blocked")
		}
		statusFilter = status
	}

	searchFilter := ""
	if search != "" {
		searchFilter = search
	}

	total, err := s.repo.PgRepo.Repo.CountDrivers(ctx, pg.CountDriversParams{
		Column1: statusFilter,
		Column2: searchFilter,
	})
	if err != nil {
		return model.ListDriversResponse{}, fmt.Errorf("failed to count drivers: %w", err)
	}

	drivers, err := s.repo.PgRepo.Repo.GetDrivers(ctx, pg.GetDriversParams{
		Column1: statusFilter,
		Column2: searchFilter,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return model.ListDriversResponse{}, fmt.Errorf("failed to get drivers: %w", err)
	}

	response := model.ListDriversResponse{
		Data: make([]model.DriverResponse, 0, len(drivers)),
		Meta: model.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}

	for _, driver := range drivers {
		response.Data = append(response.Data, s.toDriverResponse(driver))
	}

	return response, nil
}

func (s *DriverService) UpdateDriver(ctx context.Context, id string, req model.UpdateDriverRequest) (model.DriverResponse, error) {
	driverID, err := uuid.Parse(id)
	if err != nil {
		return model.DriverResponse{}, NewAppError("INVALID_ID", "invalid driver id format")
	}

	existing, err := s.repo.PgRepo.Repo.GetDriverByID(ctx, driverID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DriverResponse{}, NewAppError("DRIVER_NOT_FOUND", fmt.Sprintf("driver with id %s not found", id))
		}
		return model.DriverResponse{}, fmt.Errorf("failed to get driver: %w", err)
	}

	if err := s.validateUpdateRequest(req); err != nil {
		return model.DriverResponse{}, err
	}

	phone := existing.Phone
	if req.Phone != nil {
		phone = *req.Phone
	}

	license := existing.LicenseNumber
	if req.LicenseNumber != nil {
		license = *req.LicenseNumber
	}

	carPlate := existing.CarPlate
	if req.CarPlate != nil {
		carPlate = *req.CarPlate
	}

	if err := s.checkUniqueFields(ctx, id, phone, license, carPlate); err != nil {
		return model.DriverResponse{}, err
	}

	fullName := existing.FullName
	if req.FullName != nil {
		fullName = *req.FullName
	}

	carModel := existing.CarModel
	if req.CarModel != nil {
		carModel = *req.CarModel
	}

	driver, err := s.repo.PgRepo.Repo.UpdateDriver(ctx, pg.UpdateDriverParams{
		ID:            driverID,
		FullName:      fullName,
		Phone:         phone,
		LicenseNumber: license,
		CarModel:      carModel,
		CarPlate:      carPlate,
	})
	if err != nil {
		return model.DriverResponse{}, fmt.Errorf("failed to update driver: %w", err)
	}

	return s.toDriverResponse(driver), nil
}

func (s *DriverService) UpdateDriverStatus(ctx context.Context, id string, req model.UpdateStatusRequest) (model.DriverResponse, error) {
	driverID, err := uuid.Parse(id)
	if err != nil {
		return model.DriverResponse{}, NewAppError("INVALID_ID", "invalid driver id format")
	}

	if !isValidStatus(string(req.Status)) {
		return model.DriverResponse{}, NewAppError("INVALID_STATUS", "status must be one of: active, inactive, blocked")
	}

	driver, err := s.repo.PgRepo.Repo.UpdateDriverStatus(ctx, pg.UpdateDriverStatusParams{
		ID:     driverID,
		Status: string(req.Status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DriverResponse{}, NewAppError("DRIVER_NOT_FOUND", fmt.Sprintf("driver with id %s not found", id))
		}
		return model.DriverResponse{}, fmt.Errorf("failed to update driver status: %w", err)
	}

	return s.toDriverResponse(driver), nil
}

func (s *DriverService) DeleteDriver(ctx context.Context, id string) error {
	driverID, err := uuid.Parse(id)
	if err != nil {
		return NewAppError("INVALID_ID", "invalid driver id format")
	}

	_, err = s.repo.PgRepo.Repo.GetDriverByID(ctx, driverID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return NewAppError("DRIVER_NOT_FOUND", fmt.Sprintf("driver with id %s not found", id))
		}
		return fmt.Errorf("failed to get driver: %w", err)
	}

	if err := s.repo.PgRepo.Repo.SoftDeleteDriver(ctx, driverID); err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	return nil
}

func (s *DriverService) validateCreateRequest(req model.CreateDriverRequest) error {
	if len(req.FullName) < 3 || len(req.FullName) > 100 {
		return NewAppError("VALIDATION_ERROR", "full_name must be between 3 and 100 characters")
	}

	if !isValidPhone(req.Phone) {
		return NewAppError("VALIDATION_ERROR", "phone must be in format +998XXXXXXXXX")
	}

	if strings.TrimSpace(req.LicenseNumber) == "" {
		return NewAppError("VALIDATION_ERROR", "license_number is required")
	}

	if strings.TrimSpace(req.CarModel) == "" {
		return NewAppError("VALIDATION_ERROR", "car_model is required")
	}

	if strings.TrimSpace(req.CarPlate) == "" {
		return NewAppError("VALIDATION_ERROR", "car_plate is required")
	}

	if req.Status != "" && !isValidStatus(string(req.Status)) {
		return NewAppError("VALIDATION_ERROR", "status must be one of: active, inactive, blocked")
	}

	return nil
}

func (s *DriverService) validateUpdateRequest(req model.UpdateDriverRequest) error {
	if req.FullName != nil {
		if len(*req.FullName) < 3 || len(*req.FullName) > 100 {
			return NewAppError("VALIDATION_ERROR", "full_name must be between 3 and 100 characters")
		}
	}

	if req.Phone != nil {
		if !isValidPhone(*req.Phone) {
			return NewAppError("VALIDATION_ERROR", "phone must be in format +998XXXXXXXXX")
		}
	}

	if req.LicenseNumber != nil {
		if strings.TrimSpace(*req.LicenseNumber) == "" {
			return NewAppError("VALIDATION_ERROR", "license_number cannot be empty")
		}
	}

	if req.CarModel != nil {
		if strings.TrimSpace(*req.CarModel) == "" {
			return NewAppError("VALIDATION_ERROR", "car_model cannot be empty")
		}
	}

	if req.CarPlate != nil {
		if strings.TrimSpace(*req.CarPlate) == "" {
			return NewAppError("VALIDATION_ERROR", "car_plate cannot be empty")
		}
	}

	return nil
}

func (s *DriverService) checkUniqueFields(ctx context.Context, id, phone, license, carPlate string) error {
	var driverID uuid.UUID
	if id != "" {
		parsed, err := uuid.Parse(id)
		if err == nil {
			driverID = parsed
		}
	}

	phoneExists, err := s.repo.PgRepo.Repo.CheckPhoneExists(ctx, pg.CheckPhoneExistsParams{
		Phone:   phone,
		Column2: driverID,
	})
	if err != nil {
		return fmt.Errorf("failed to check phone: %w", err)
	}
	if phoneExists {
		return NewAppError("PHONE_CONFLICT", "phone number already exists")
	}

	licenseExists, err := s.repo.PgRepo.Repo.CheckLicenseExists(ctx, pg.CheckLicenseExistsParams{
		LicenseNumber: license,
		Column2:       driverID,
	})
	if err != nil {
		return fmt.Errorf("failed to check license: %w", err)
	}
	if licenseExists {
		return NewAppError("LICENSE_CONFLICT", "license number already exists")
	}

	plateExists, err := s.repo.PgRepo.Repo.CheckCarPlateExists(ctx, pg.CheckCarPlateExistsParams{
		CarPlate: carPlate,
		Column2:  driverID,
	})
	if err != nil {
		return fmt.Errorf("failed to check car plate: %w", err)
	}
	if plateExists {
		return NewAppError("CAR_PLATE_CONFLICT", "car plate already exists")
	}

	return nil
}

func (s *DriverService) toDriverResponse(driver pg.Driver) model.DriverResponse {
	return model.DriverResponse{
		ID:            driver.ID.String(),
		FullName:      driver.FullName,
		Phone:         driver.Phone,
		LicenseNumber: driver.LicenseNumber,
		CarModel:      driver.CarModel,
		CarPlate:      driver.CarPlate,
		Status:        model.DriverStatus(driver.Status),
		CreatedAt:     driver.CreatedAt.Time,
		UpdatedAt:     driver.UpdatedAt.Time,
	}
}

func isValidPhone(phone string) bool {
	match, _ := regexp.MatchString(`^\+998[0-9]{9}$`, phone)
	return match
}

func isValidStatus(status string) bool {
	return status == "active" || status == "inactive" || status == "blocked"
}

type AppError struct {
	Code    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewAppError(code, message string) AppError {
	return AppError{
		Code:    code,
		Message: message,
	}
}

func (s *DriverService) GetActiveDriversStats(ctx context.Context, status string) (model.ActiveDriversStatsResponse, error) {
	status = strings.TrimSpace(status)
	if status == "" {
		status = string(model.StatusActive)
	}
	if !isValidStatus(status) {
		return model.ActiveDriversStatsResponse{}, NewAppError("INVALID_STATUS", "status must be one of: active, inactive, blocked")
	}

	count, err := s.repo.PgRepo.Repo.CountDrivers(ctx, pg.CountDriversParams{
		Column1: status,
		Column2: "",
	})
	if err != nil {
		return model.ActiveDriversStatsResponse{}, fmt.Errorf("failed to get drivers count by status: %w", err)
	}

	countInt := int(count)
	return model.ActiveDriversStatsResponse{
		Status: model.DriverStatus(status),
		Count:  countInt,
	}, nil
}
