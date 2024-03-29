// Package client provides an authsvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	"io"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"remote-storage/server/authsvc"
)

// New returns a service that's load-balanced over instances of authsvc found
// in the provided Consul server. The mechanism of looking up authsvc
// instances in Consul is hard-coded into the client.
func New(consulAddr string, logger log.Logger) (authsvc.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	var (
		consulService = "auth-service"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints authsvc.Endpoints
	)
	{
		factory := factoryFor(authsvc.MakeLoginEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.LoginEndpoint = retry
	}
	{
		factory := factoryFor(authsvc.MakeRefreshTokenEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.RefreshTokenEndpoint = retry
	}
	{
		factory := factoryFor(authsvc.MakeValidateTokenEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.ValidateTokenEndpoint = retry
	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(authsvc.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := authsvc.MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
