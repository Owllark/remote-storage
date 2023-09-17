// Package client provides a storagesvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	"io"
	"net/http"
	"remote-storage/server/storagesvc"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

type Service interface {
	storagesvc.Service
	// it seems to me that AddCookie method breaks the abstraction
	// but currently i didn't find a better way to set nescessary cookies
	// (authentication cookies firstly that are returning by Authentication Service)
	AddCookie(cookie *http.Cookie)
}

// New returns a service that's load-balanced over instances of storagesvc found
// in the provided Consul server. The mechanism of looking up storagesvc
// instances in Consul is hard-coded into the client.
func New(consulAddr string, logger log.Logger) (Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	var (
		consulService = "file-system-service"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints storagesvc.Endpoints
	)
	{
		factory := factoryFor(storagesvc.MakeGetStateEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetStateEndpoint = retry
	}
	{
		factory := factoryFor(storagesvc.MakeMkDirEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.MkDirEndpoint = retry
	}
	{
		factory := factoryFor(storagesvc.MakeRenameEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.RenameEndpoint = retry
	}
	{
		factory := factoryFor(storagesvc.MakeMoveEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.MoveEndpoint = retry
	}
	{
		factory := factoryFor(storagesvc.MakeCopyEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.CopyEndpoint = retry
	}
	{
		factory := factoryFor(storagesvc.MakeDeleteEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.DeleteEndpoint = retry
	}
	{
		factory := factoryFor(storagesvc.MakeDownloadEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.DownloadEndpoint = retry
	}
	{
		factory := factoryFor(storagesvc.MakeUploadEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.UploadEndpoint = retry
	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(storagesvc.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := storagesvc.MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
