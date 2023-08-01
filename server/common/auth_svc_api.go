package common

type LoginRequest struct {
	Login          string `json:"login,omitempty"`
	HashedPassword string `json:"hashed_password,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"err,omitempty"`
}

type RefreshTokenRequest struct {
	Token string `json:"token,omitempty"`
}

type RefreshTokenResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"err,omitempty"`
}

type ValidateTokenRequest struct {
	Token string `json:"token,omitempty"`
}

type ValidateTokenResponse struct {
	Inf   UserInf `json:"inf"`
	Error string  `json:"err,omitempty"`
}
