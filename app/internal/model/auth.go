package model

type User struct {
	ID          string  `json:"id"`
	FullName    *string `json:"fullName,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

// LoginRequest represents login request body
type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber" example:"+998901234567"`
	Password    string `json:"password" example:"ASDF1234"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	AccessToken  string       `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         UserResponse `json:"user"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
}

// RegisterRequest represents registration request body
type RegisterRequest struct {
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password" validate:"required,min=8,containsLettersAndNumbers" example:"ASDF1234"`
}

// UserResponse represents a safe subset of user data returned to clients
type UserResponse struct {
	ID          string  `json:"id"`
	FullName    *string `json:"fullName,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type LogoutRequest struct {
	AccessToken string `json:"accessToken"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type LogoutResponse struct {
	Message string `json:"message"`
}
