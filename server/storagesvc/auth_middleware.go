package storagesvc

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"remote-storage/server/authsvc"
	"remote-storage/server/common"
)

type ctxUserInfKey struct{}

func putUserInfInCtx(ctx context.Context, r common.UserInf) context.Context {
	return context.WithValue(ctx, ctxUserInfKey{}, r)
}

func AuthMiddleware(authSvc authsvc.Service, next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// Extracting the request from the context
		req := ctx.Value(ctxRequestKey{}).(*http.Request)
		// Accessing the cookies from the request
		tokenCookie, err := req.Cookie("token")
		if err != nil {
			return nil, ErrAuthFailed
		}

		userInf, err := authSvc.ValidateToken(ctx, tokenCookie.Value)
		if err != nil {
			return nil, err
		}

		putUserInfInCtx(ctx, userInf)

		// Call the next endpoint in the chain

		return next(ctx, request)
	}
}
