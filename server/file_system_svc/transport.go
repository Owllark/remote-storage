package file_system_svc

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type ctxRequestKey struct{}

func putRequestInCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, ctxRequestKey{}, r)
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

func ApplyTransportMiddleware(endpoints *Endpoints, mw TransportMiddleware) Endpoints {
	v := reflect.ValueOf(endpoints).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Type() == reflect.TypeOf((*endpoint.Endpoint)(nil)).Elem() {
			endpointFunc := field.Interface().(endpoint.Endpoint)
			field.Set(reflect.ValueOf(mw(endpointFunc)))
		}
	}
	return *endpoints
}

// MakeHttpHandler mounts all service endpoints into a http.Handler.
// Useful in a profilesvc server.
func MakeHttpHandler(s FileSystemService) http.Handler {

	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	e = ApplyTransportMiddleware(&e, AuthMiddleware)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(putRequestInCtx),
	}

	r.Methods("GET").Path("/filesystem/test").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(917)
	})

	r.Methods("GET").Path("/filesystem/state").Handler(kithttp.NewServer(
		e.GetStateEndpoint,
		decodeGetStateRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/filesystem/mkdir").Handler(kithttp.NewServer(
		e.MkDirEndpoint,
		decodeMkDirRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/filesystem/rename").Handler(kithttp.NewServer(
		e.RenameEndpoint,
		decodeRenameRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/filesystem/copy").Handler(kithttp.NewServer(
		e.CopyEndpoint,
		decodeCopyRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/filesystem/move").Handler(kithttp.NewServer(
		e.MoveEndpoint,
		decodeMoveRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("filesystem/delete").Handler(kithttp.NewServer(
		e.DeleteEndpoint,
		decodeDeleteRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeGetStateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getStateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeMkDirRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request mkDirRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeRenameRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request renameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCopyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request copyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeMoveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request moveRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request deleteRequest
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
	case ErrAuthFailed.Error():
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
