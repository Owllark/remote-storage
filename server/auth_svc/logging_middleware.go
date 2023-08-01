package auth_svc

import (
	"github.com/go-kit/kit/log"
	"server/common"
	"time"
)

type Middleware func(AuthService) AuthService

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next AuthService) AuthService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   AuthService
	logger log.Logger
}

func (mw loggingMiddleware) Login(login string, password string) (token string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Login", "login", login, "password", password, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Login(login, password)
}

func (mw loggingMiddleware) ValidateToken(token string) (inf common.UserInf, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "ValidateToken", "token", token, time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.ValidateToken(token)
}

func (mw loggingMiddleware) RefreshToken(token string) (newToken string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "RefrashToken", "token", token, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.RefreshToken(token)
}
