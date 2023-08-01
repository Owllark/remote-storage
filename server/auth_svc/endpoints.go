package auth_svc

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"server/common"
)

type Endpoints struct {
	LoginEndpoint         endpoint.Endpoint
	RefreshTokenEndpoint  endpoint.Endpoint
	ValidateTokenEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the FileSystemService.
func MakeServerEndpoints(s AuthService) Endpoints {
	return Endpoints{
		LoginEndpoint:         makeLoginEndpoint(s),
		RefreshTokenEndpoint:  makeRefreshTokenEndpoint(s),
		ValidateTokenEndpoint: makeValidateTokenEndpoint(s),
	}
}

func makeLoginEndpoint(svc AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		token, err := svc.Login(req.Login, req.HashedPassword)
		if err != nil {
			return loginResponse{token, err.Error()}, nil
		}

		return loginResponse{token, ""}, nil
	}
}

func makeRefreshTokenEndpoint(svc AuthService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(refreshTokenRequest)
		resp, err := svc.RefreshToken(req.Token)
		if err != nil {
			return refreshTokenResponse{resp, err.Error()}, nil
		}
		return refreshTokenResponse{resp, ""}, nil
	}
}

func makeValidateTokenEndpoint(svc AuthService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(validateTokenRequest)
		resp, err := svc.ValidateToken(req.Token)
		if err != nil {
			return validateTokenResponse{resp, err.Error()}, nil
		}
		return validateTokenResponse{resp, ""}, nil
	}
}

type loginRequest struct {
	Login          string `json:"login,omitempty"`
	HashedPassword string `json:"hashed_password,omitempty"`
}

type loginResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"err,omitempty"`
}

func (r loginResponse) error() error {
	return errors.New(r.Error)
}

type refreshTokenRequest struct {
	Token string `json:"token,omitempty"`
}

type refreshTokenResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"err,omitempty"`
}

func (r refreshTokenResponse) error() error {
	return errors.New(r.Error)
}

type validateTokenRequest struct {
	Token string `json:"token,omitempty"`
}

type validateTokenResponse struct {
	Inf   common.UserInf `json:"inf"`
	Error string         `json:"err,omitempty"`
}

func (r validateTokenResponse) error() error {
	return errors.New(r.Error)
}
