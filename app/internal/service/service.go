package service

import (
	"context"

	"gitlab.yurtal.tech/company/blitz/back/internal/config"
	"gitlab.yurtal.tech/company/blitz/back/internal/model"
	"gitlab.yurtal.tech/company/blitz/back/internal/repository"
)

type AuthI interface {
	Register(ctx context.Context, req model.RegisterRequest, jwtCfg *config.JwtConfig) (model.LoginResponse, error)
	Login(ctx context.Context, req model.LoginRequest, jwtCfg *config.JwtConfig) (model.LoginResponse, error)
	Refresh(ctx context.Context, req model.RefreshRequest, jwtCfg *config.JwtConfig) (model.RefreshResponse, error)
}

type DriverI interface {
	CreateDriver(ctx context.Context, req model.CreateDriverRequest) (model.DriverResponse, error)
	GetDriver(ctx context.Context, id string) (model.DriverResponse, error)
	ListDrivers(ctx context.Context, page, limit int, status, search string) (model.ListDriversResponse, error)
	UpdateDriver(ctx context.Context, id string, req model.UpdateDriverRequest) (model.DriverResponse, error)
	UpdateDriverStatus(ctx context.Context, id string, req model.UpdateStatusRequest) (model.DriverResponse, error)
	DeleteDriver(ctx context.Context, id string) error
	LogActivity(ctx context.Context, driverID string, req model.LogActivityRequest) (model.ActivityLogResponse, error)
	GetActivityLog(ctx context.Context, driverID string, page, limit int) ([]model.ActivityLogResponse, error)
	GetWorkingHours(ctx context.Context, driverID, startDate, endDate string) (model.WorkingHoursResponse, error)
	GetActiveDriversStats(ctx context.Context) (model.ActiveDriversStatsResponse, error)
}

type I interface {
	Auth() AuthI
	Driver() DriverI
}

type Service struct {
	auth   AuthI
	driver DriverI
}

func New(cfg *config.Config, repo *repository.Repository) *Service {
	return &Service{
		auth:   NewAuthS(cfg, repo),
		driver: NewDriverService(cfg, repo),
	}
}

func (s *Service) Auth() AuthI {
	return s.auth
}

func (s *Service) Driver() DriverI {
	return s.driver
}
