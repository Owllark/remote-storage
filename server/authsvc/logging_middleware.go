package authsvc

import (
	"context"
	"github.com/go-kit/kit/log"
	"server/common"
	"time"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) Login(ctx context.Context, login string, password string) (token string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Login", "login", password, "password", login, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Login(ctx, login, password)
}

func (mw loggingMiddleware) ValidateToken(ctx context.Context, tokenStr string) (inf common.UserInf, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "ValidateToken", "token", tokenStr, time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.ValidateToken(ctx, tokenStr)
}

func (mw loggingMiddleware) RefreshToken(ctx context.Context, tokenStr string) (newToken string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "RefrashToken", "token", tokenStr, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.RefreshToken(ctx, tokenStr)
}
