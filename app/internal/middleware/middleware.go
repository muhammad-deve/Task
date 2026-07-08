// app/internal/middleware/middleware.go
package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.yurtal.tech/company/blitz/back/internal/config"
	"gitlab.yurtal.tech/company/blitz/back/internal/model"
	"gitlab.yurtal.tech/company/blitz/back/pkg/utils"
)

func SetupMiddleware(e *echo.Echo, cfg *config.Config) {
	// Request ID
	e.Use(middleware.RequestID())

	// Logger
	e.Use(middleware.Logger())

	// Recover from panics
	e.Use(middleware.Recover())

	// CORS - accept requests from any origin (mirrors the request Origin back)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Authorization", "Content-Disposition"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	// Timeout
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(cfg.Server.CtxDefaultTimeout) * time.Second,
	}))

	// Rate limiter
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
}

// LoginRateLimiter limits login attempts per IP to mitigate brute-force attacks
func LoginRateLimiter() echo.MiddlewareFunc {
	store := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      5,
		Burst:     5,
		ExpiresIn: time.Minute,
	})

	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: store,
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusTooManyRequests, model.ErrorResponse{Message: "too many login attempts"})
		},
	})
}

// CheckAuth validates a JWT access token and sets user_id in context.
// It accepts Authorization: Bearer <token> and Authorization: <token>.
func CheckAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var accessToken string

			authHeader := strings.TrimSpace(c.Request().Header.Get("Authorization"))
			fields := strings.Fields(authHeader)
			if len(fields) == 1 {
				accessToken = fields[0]
			} else if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
				accessToken = fields[1]
			}

			if accessToken == "" {
				return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "you are not logged in"})
			}

			sub, err := utils.ValidateJWT(accessToken, cfg.Jwt.SecretKey)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: err.Error()})
			}

			c.Set("user_id", fmt.Sprint(sub))
			return next(c)
		}
	}
}

// ValidateLoginInput binds and validates login request body before reaching handler
func ValidateLoginInput(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
		}
		if req.PhoneNumber == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "phone number and password are required"})
		}

		c.Set("loginBody", req)
		return next(c)
	}
}

// ValidateRegisterInput binds and validates registration request body before reaching handler
func ValidateRegisterInput(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.RegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
		}
		if req.PhoneNumber == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "phone number and password are required"})
		}

		// Validate password: at least 8 characters with letters and numbers
		if !utils.ValidatePassword(req.Password) {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "password must be at least 8 characters and contain both letters and numbers"})
		}

		c.Set("registerBody", req)
		return next(c)
	}

}

// ValidateRefreshInput binds and validates refresh token request body before reaching handler
func ValidateRefreshInput(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.RefreshRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "invalid request body"})
		}
		if req.RefreshToken == "" {
			return c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "refresh token is required"})
		}

		c.Set("refreshBody", req)
		return next(c)
	}
}
