package file_system_svc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"io"
	"net/http"
	"server/common"
)

type ctxUserInfKey struct{}

func putUserInfInCtx(ctx context.Context, r common.UserInf) context.Context {
	return context.WithValue(ctx, ctxUserInfKey{}, r)
}

const auth_svc_url = "http://localhost:666/authentication/validate"

func AuthMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// Extracting the request from the context
		req := ctx.Value(ctxRequestKey{}).(*http.Request)
		// Accessing the cookies from the request
		tokenCookie, err := req.Cookie("token")
		if err != nil {
			return nil, ErrAuthFailed
		}

		authRequestBody := common.ValidateTokenRequest{Token: tokenCookie.Value}
		authRequestBodyJson, err := json.Marshal(authRequestBody)
		if err != nil {
			return nil, ErrUnknownError
		}
		resp, err := http.Post(auth_svc_url, "application/json", bytes.NewReader(authRequestBodyJson))
		body, _ := io.ReadAll(resp.Body)
		var authResponse common.ValidateTokenResponse
		err = json.Unmarshal(body, &authResponse)
		if err != nil {
			return nil, ErrUnknownError
		}
		if authResponse.Error != "" {
			return nil, errors.New(authResponse.Error)
		}

		putUserInfInCtx(ctx, authResponse.Inf)

		// Call the next endpoint in the chain

		return next(ctx, request)
	}
}
