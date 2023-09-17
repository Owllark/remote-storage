package authsvc

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/url"
	"remote-storage/server/common"
	"strings"
	"time"
)

type Endpoints struct {
	LoginEndpoint         endpoint.Endpoint
	RefreshTokenEndpoint  endpoint.Endpoint
	ValidateTokenEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the Authentication Service.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		LoginEndpoint:         MakeLoginEndpoint(s),
		RefreshTokenEndpoint:  MakeRefreshTokenEndpoint(s),
		ValidateTokenEndpoint: MakeValidateTokenEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in the authsvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return Endpoints{
		LoginEndpoint:         httptransport.NewClient("POST", tgt, encodeLoginRequest, decodeLoginResponse, options...).Endpoint(),
		RefreshTokenEndpoint:  httptransport.NewClient("POST", tgt, encodeRefreshTokenRequest, decodeRefreshTokenResponse, options...).Endpoint(),
		ValidateTokenEndpoint: httptransport.NewClient("POST", tgt, encodeValidateTokenRequest, decodeValidateTokenResponse, options...).Endpoint(),
	}, nil
}

// Login implements Service. Primarily useful in a client.
func (e Endpoints) Login(ctx context.Context, login string, password string) (AuthCookie, error) {
	request := loginRequest{
		Login:          login,
		HashedPassword: password,
	}
	response, err := e.LoginEndpoint(ctx, request)
	if err != nil {
		return AuthCookie{}, err
	}
	resp := response.(loginResponse)
	return resp.AuthCookie, errors.New(resp.Error)
}

// RefreshToken implements Service. Primarily useful in a client.
func (e Endpoints) RefreshToken(ctx context.Context, tokenStr string) (AuthCookie, error) {
	request := refreshTokenRequest{
		Token: tokenStr,
	}
	response, err := e.RefreshTokenEndpoint(ctx, request)
	if err != nil {
		return AuthCookie{}, err
	}
	resp := response.(refreshTokenResponse)
	return resp.AuthCookie, errors.New(resp.Error)
}

// ValidateToken implements Service. Primarily useful in a client.
func (e Endpoints) ValidateToken(ctx context.Context, tokenStr string) (common.UserInf, error) {
	request := validateTokenRequest{
		Token: tokenStr,
	}
	response, err := e.ValidateTokenEndpoint(ctx, request)
	if err != nil {
		return common.UserInf{}, err
	}
	resp := response.(validateTokenResponse)
	return resp.Inf, errors.New(resp.Error)
}

func MakeLoginEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		token, err := svc.Login(ctx, req.Login, req.HashedPassword)
		if err != nil {
			return loginResponse{token, err.Error()}, nil
		}

		return loginResponse{token, ""}, nil
	}
}

func MakeRefreshTokenEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(refreshTokenRequest)
		resp, err := svc.RefreshToken(ctx, req.Token)
		if err != nil {
			return refreshTokenResponse{resp, err.Error()}, nil
		}
		return refreshTokenResponse{resp, ""}, nil
	}
}

func MakeValidateTokenEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateTokenRequest)
		resp, err := svc.ValidateToken(ctx, req.Token)
		if err != nil {
			return validateTokenResponse{resp, err.Error()}, nil
		}
		return validateTokenResponse{resp, ""}, nil
	}
}

type AuthCookie struct {
	Name    string    `json:"name,omitempty"`
	Value   string    `json:"value,omitempty"`
	Expires time.Time `json:"expires"`
}

type loginRequest struct {
	Login          string `json:"login,omitempty"`
	HashedPassword string `json:"hashed_password,omitempty"`
}

type loginResponse struct {
	AuthCookie AuthCookie `json:"auth_cookie,omitempty"`
	Error      string     `json:"error,omitempty"`
}

func (r loginResponse) error() error {
	return errors.New(r.Error)
}

type refreshTokenRequest struct {
	Token string `json:"token,omitempty"`
}

type refreshTokenResponse struct {
	AuthCookie AuthCookie `json:"auth_cookie,omitempty"`
	Error      string     `json:"error,omitempty"`
}

func (r refreshTokenResponse) error() error {
	return errors.New(r.Error)
}

type validateTokenRequest struct {
	Token string `json:"token,omitempty"`
}

type validateTokenResponse struct {
	Inf   common.UserInf `json:"inf"`
	Error string         `json:"error,omitempty"`
}

func (r validateTokenResponse) error() error {
	return errors.New(r.Error)
}
