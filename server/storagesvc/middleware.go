package storagesvc

import (
	"github.com/go-kit/kit/endpoint"
	"remote-storage/server/authsvc"
)

type Middleware func(Service) Service
type TransportMiddleware func(endpoint.Endpoint) endpoint.Endpoint
type TransportAuthMiddleware func(authsvc.Service, endpoint.Endpoint) endpoint.Endpoint
