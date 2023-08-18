package file_system_svc

import (
	"github.com/go-kit/kit/endpoint"
	"server/authsvc"
)

type Middleware func(FileSystemService) FileSystemService
type TransportMiddleware func(endpoint.Endpoint) endpoint.Endpoint
type TransportAuthMiddleware func(authsvc.Service, endpoint.Endpoint) endpoint.Endpoint
