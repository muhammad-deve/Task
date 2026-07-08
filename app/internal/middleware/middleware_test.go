package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.yurtal.tech/company/blitz/back/internal/config"
	"gitlab.yurtal.tech/company/blitz/back/pkg/utils"
)

func TestCheckAuthAcceptsBearerAndRawAuthorization(t *testing.T) {
	cfg := &config.Config{}
	cfg.Jwt.SecretKey = "test-secret"

	token, err := utils.CreateJWT(time.Hour, "user-1", cfg.Jwt.SecretKey)
	if err != nil {
		t.Fatalf("CreateJWT() error = %v", err)
	}

	tests := []struct {
		name          string
		authorization string
		wantStatus    int
	}{
		{
			name:          "bearer token",
			authorization: "Bearer " + token,
			wantStatus:    http.StatusOK,
		},
		{
			name:          "raw token",
			authorization: token,
			wantStatus:    http.StatusOK,
		},
		{
			name:       "missing token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := CheckAuth(cfg)(func(c echo.Context) error {
				if c.Get("user_id") == "" {
					t.Fatal("expected user_id to be set")
				}
				return c.NoContent(http.StatusOK)
			})

			if err := handler(c); err != nil {
				t.Fatalf("handler() error = %v", err)
			}

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
