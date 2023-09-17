package storagesvc

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"reflect"
	"remote-storage/server/authsvc"
)

const chunkSize = 64 * 1

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

type readCloserContainer interface {
	ReadCloser() io.ReadCloser
}

func ApplyAuthMiddleware(endpoints *Endpoints, mw TransportAuthMiddleware, svc authsvc.Service) Endpoints {
	v := reflect.ValueOf(endpoints).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Type() == reflect.TypeOf((*endpoint.Endpoint)(nil)).Elem() {
			endpointFunc := field.Interface().(endpoint.Endpoint)
			field.Set(reflect.ValueOf(mw(svc, endpointFunc)))
		}
	}
	return *endpoints
}

// MakeHttpHandler mounts all service endpoints into a http.Handler.
// Useful in a storagesvc server.
func MakeHttpHandler(s Service) http.Handler {

	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	e = ApplyAuthMiddleware(&e, AuthMiddleware, s.getAuthSvc())
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(putRequestInCtx),
	}

	r.Methods("GET").Path("/filesystem/state").Handler(kithttp.NewServer(
		e.GetStateEndpoint,
		decodeGetStateRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/filesystem/download").Handler(kithttp.NewServer(
		e.DownloadEndpoint,
		decodeDownloadRequest,
		encodeChunkedResponse,
		options...,
	))
	r.Methods("POST").Path("/filesystem/upload").Handler(kithttp.NewServer(
		e.UploadEndpoint,
		decodeUploadRequest,
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

func decodeDownloadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request downloadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUploadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uploadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorer)
	if !ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeChunkedResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorer)
	if !ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		encodeError(ctx, e.error(), w)
		return nil
	}
	r, ok := response.(readCloserContainer)
	if !ok {
		encodeError(ctx, ErrUnknownError, w)
		return nil
	}
	srcBuffer := r.ReadCloser()
	buffer := make([]byte, chunkSize)
	for {
		// Read the next chunk from the file
		n, err := srcBuffer.Read(buffer)
		if err == io.EOF {
			// Reached the end of the file, break the loop
			break
		}
		if err != nil {
			// Error reading the file, handle it
			encodeError(ctx, ErrUnknownError, w)
			return nil
		}

		// Send the current chunk to the response writer
		_, err = w.Write(buffer[:n])
		if err != nil {
			// Error sending the chunk, handle it
			encodeError(ctx, ErrUnknownError, w)
			return nil
		}

		// Flush the response writer to ensure the chunk is sent immediately
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
	return nil
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// storagesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(&buf)
	cookies := ctx.Value(contextKeyRequestCookie).(Cookies)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return nil
}

func encodeGetStateRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	// r.Methods("GET").Path("/filesystem/state")
	req.URL.Path = "/filesystem/state"
	return encodeRequest(ctx, req, reqBody)
}

func decodeGetStateResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response getStateResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeMkDirRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	// r.Methods("POST").Path("/filesystem/mkdir")
	req.URL.Path = "/filesystem/mkdir"
	return encodeRequest(ctx, req, reqBody)
}

func decodeMkDirResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response mkDirResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeRenameRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	// r.Methods("POST").Path("/filesystem/rename")
	req.URL.Path = "/filesystem/rename"
	return encodeRequest(ctx, req, reqBody)
}

func decodeRenameResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response renameResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeMoveRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	// r.Methods("POST").Path("/filesystem/move")
	req.URL.Path = "/filesystem/move"
	return encodeRequest(ctx, req, reqBody)
}

func decodeMoveResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response moveResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeCopyRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	// r.Methods("POST").Path("/filesystem/copy")
	req.URL.Path = "/filesystem/copy"
	return encodeRequest(ctx, req, reqBody)
}

func decodeCopyResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response moveResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeDeleteRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	// r.Methods("POST").Path("/filesystem/delete")
	req.URL.Path = "/filesystem/delete"
	return encodeRequest(ctx, req, reqBody)
}

func decodeDeleteResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response moveResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeDownloadRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	// r.Methods("POST").Path("/filesystem/download")
	req.URL.Path = "/filesystem/download"
	return encodeRequest(ctx, req, reqBody)
}

func decodeDownloadResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response moveResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeUploadRequest(ctx context.Context, req *http.Request, reqBody interface{}) error {
	// r.Methods("POST").Path("/filesystem/upload")
	req.URL.Path = "/filesystem/upload"
	return encodeRequest(ctx, req, reqBody)
}

func decodeUdploadResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response moveResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
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
