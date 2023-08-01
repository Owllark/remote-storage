package file_system_svc

import "github.com/go-kit/kit/endpoint"

type Middleware func(FileSystemService) FileSystemService
type TransportMiddleware func(endpoint.Endpoint) endpoint.Endpoint
