package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"server/file_system_svc"
	"time"
)

//
//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"log"
//	"net/http"
//	"strings"
//
//	"github.com/go-kit/kit/endpoint"
//	httptransport "github.com/go-kit/kit/transport/http"
//)
//
//// StringService provides operations on strings.
//type StringService interface {
//	Uppercase(string) (string, error)
//	Count(string) int
//}
//
//// stringService is a concrete implementation of StringService
//type stringService struct{}
//
//func (stringService) Uppercase(s string) (string, error) {
//	if s == "" {
//		return "", ErrEmpty
//	}
//	return strings.ToUpper(s), nil
//}
//
//func (stringService) Count(s string) int {
//	return len(s)
//}
//
//// ErrEmpty is returned when an input string is empty.
//var ErrEmpty = errors.New("empty string")
//
//// For each method, we define request and response structs
//type uppercaseRequest struct {
//	S string `json:"s"`
//}
//
//type uppercaseResponse struct {
//	V   string `json:"v"`
//	Error string `json:"err,omitempty"` // errors don't define JSON marshaling
//}
//
//type countRequest struct {
//	S string `json:"s"`
//}
//
//type countResponse struct {
//	V int `json:"v"`
//}
//
//// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
//func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
//	return func(_ context.Context, request interface{}) (interface{}, error) {
//		req := request.(uppercaseRequest)
//		v, err := svc.Uppercase(req.S)
//		if err != nil {
//			return uppercaseResponse{v, err.Error()}, nil
//		}
//		return uppercaseResponse{v, ""}, nil
//	}
//}
//
//func makeCountEndpoint(svc StringService) endpoint.Endpoint {
//	return func(_ context.Context, request interface{}) (interface{}, error) {
//		req := request.(countRequest)
//		v := svc.Count(req.S)
//		return countResponse{v}, nil
//	}
//}
//
//// Transports expose the service to the network. In this first example we utilize JSON over HTTP.
//func main() {
//	svc := stringService{}
//
//	uppercaseHandler := httptransport.NewServer(
//		makeUppercaseEndpoint(svc),
//		decodeUppercaseRequest,
//		encodeResponse,
//	)
//
//	countHandler := httptransport.NewServer(
//		makeCountEndpoint(svc),
//		decodeCountRequest,
//		encodeResponse,
//	)
//
//	http.Handle("/uppercase", uppercaseHandler)
//	http.Handle("/count", countHandler)
//	log.Fatal(http.ListenAndServe(":777", nil))
//}
//
//func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	var request uppercaseRequest
//	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
//		return nil, err
//	}
//	return request, nil
//}
//
//func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	var request countRequest
//	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
//		return nil, err
//	}
//	return request, nil
//}
//
//func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
//	return json.NewEncoder(w).Encode(response)
//}

func main() {
	logFile, err := os.OpenFile("file_system_svc\\log\\log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(logFile)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s file_system_svc.FileSystemService
	{
		s = file_system_svc.NewFileSystemService()
		s = file_system_svc.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = file_system_svc.MakeHttpHandler(s)
	}

	http.ListenAndServe("localhost:777", h)
	logger.Log("ServiceSessionEnd", time.Now())
}
