package service

import (
	"testing"

	"gitlab.yurtal.tech/company/blitz/back/internal/model"
)

func TestValidateCreateRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     model.CreateDriverRequest
		wantErr bool
		errCode string
	}{
		{
			name: "valid request",
			req: model.CreateDriverRequest{
				FullName:      "Akmal Karimov",
				Phone:         "+998901234567",
				LicenseNumber: "AB1234567",
				CarModel:      "Chevrolet Nexia",
				CarPlate:      "01A123BC",
			},
			wantErr: false,
		},
		{
			name: "full_name too short",
			req: model.CreateDriverRequest{
				FullName:      "AK",
				Phone:         "+998901234567",
				LicenseNumber: "AB1234567",
				CarModel:      "Chevrolet Nexia",
				CarPlate:      "01A123BC",
			},
			wantErr: true,
			errCode: "VALIDATION_ERROR",
		},
		{
			name: "full_name too long",
			req: model.CreateDriverRequest{
				FullName:      "This is a very long name that exceeds the maximum allowed length of one hundred characters and should fail",
				Phone:         "+998901234567",
				LicenseNumber: "AB1234567",
				CarModel:      "Chevrolet Nexia",
				CarPlate:      "01A123BC",
			},
			wantErr: true,
			errCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid phone format",
			req: model.CreateDriverRequest{
				FullName:      "Akmal Karimov",
				Phone:         "+998123",
				LicenseNumber: "AB1234567",
				CarModel:      "Chevrolet Nexia",
				CarPlate:      "01A123BC",
			},
			wantErr: true,
			errCode: "VALIDATION_ERROR",
		},
		{
			name: "phone without country code",
			req: model.CreateDriverRequest{
				FullName:      "Akmal Karimov",
				Phone:         "901234567",
				LicenseNumber: "AB1234567",
				CarModel:      "Chevrolet Nexia",
				CarPlate:      "01A123BC",
			},
			wantErr: true,
			errCode: "VALIDATION_ERROR",
		},
		{
			name: "empty license number",
			req: model.CreateDriverRequest{
				FullName:      "Akmal Karimov",
				Phone:         "+998901234567",
				LicenseNumber: "",
				CarModel:      "Chevrolet Nexia",
				CarPlate:      "01A123BC",
			},
			wantErr: true,
			errCode: "VALIDATION_ERROR",
		},
		{
			name: "empty car model",
			req: model.CreateDriverRequest{
				FullName:      "Akmal Karimov",
				Phone:         "+998901234567",
				LicenseNumber: "AB1234567",
				CarModel:      "",
				CarPlate:      "01A123BC",
			},
			wantErr: true,
			errCode: "VALIDATION_ERROR",
		},
		{
			name: "empty car plate",
			req: model.CreateDriverRequest{
				FullName:      "Akmal Karimov",
				Phone:         "+998901234567",
				LicenseNumber: "AB1234567",
				CarModel:      "Chevrolet Nexia",
				CarPlate:      "",
			},
			wantErr: true,
			errCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid status",
			req: model.CreateDriverRequest{
				FullName:      "Akmal Karimov",
				Phone:         "+998901234567",
				LicenseNumber: "AB1234567",
				CarModel:      "Chevrolet Nexia",
				CarPlate:      "01A123BC",
				Status:        "invalid_status",
			},
			wantErr: true,
			errCode: "VALIDATION_ERROR",
		},
	}

	s := &DriverService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.validateCreateRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				appErr, ok := err.(AppError)
				if !ok {
					t.Errorf("expected AppError, got %T", err)
					return
				}
				if appErr.Code != tt.errCode {
					t.Errorf("expected error code %s, got %s", tt.errCode, appErr.Code)
				}
			}
		})
	}
}

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{"valid phone", "+998901234567", true},
		{"valid phone 2", "+998991234567", true},
		{"missing plus", "998901234567", false},
		{"wrong country code", "+997901234567", false},
		{"too short", "+99890123456", false},
		{"too long", "+9989012345678", false},
		{"contains letters", "+998a01234567", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPhone(tt.phone); got != tt.want {
				t.Errorf("isValidPhone(%s) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"active", "active", true},
		{"inactive", "inactive", true},
		{"blocked", "blocked", true},
		{"invalid", "invalid", false},
		{"empty", "", false},
		{"uppercase", "ACTIVE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidStatus(tt.status); got != tt.want {
				t.Errorf("isValidStatus(%s) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
