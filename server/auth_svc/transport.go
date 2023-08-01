package auth_svc

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// MakeHttpHandler mounts all service endpoints into a http.Handler.
// Useful in a profilesvc server.
func MakeHttpHandler(s AuthService) http.Handler {

	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/authentication/test").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(917)
	})

	r.Methods("POST").Path("/authentication/login").Handler(httptransport.NewServer(
		e.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/authentication/validate").Handler(httptransport.NewServer(
		e.ValidateTokenEndpoint,
		decodeValidateTokenRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/authentication/refresh").Handler(httptransport.NewServer(
		e.RefreshTokenEndpoint,
		decodeRefreshTokenRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request loginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeValidateTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request validateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request refreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorer)
	if !ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err.Error() {
	case ErrNotFound.Error():
		return http.StatusNotFound
	case ErrAlreadyExists.Error():
		return http.StatusBadRequest
	case ErrWrongCredentials.Error():
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
