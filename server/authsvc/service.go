package authsvc

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"server/common"
	database "server/db"
	"time"
)

type Service interface {
	Login(ctx context.Context, login string, password string) (string, error)
	RefreshToken(ctx context.Context, tokenStr string) (string, error)
	ValidateToken(ctx context.Context, tokenStr string) (common.UserInf, error)
}

var (
	ErrUnknownError     = errors.New("unknown error")
	ErrAlreadyExists    = errors.New("already exists")
	ErrNotFound         = errors.New("not found")
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrTokenExpired     = errors.New("token expired")
)

type Claims struct {
	Username string `json:"username"`
	RootDir  string `json:"root_dir"`
	jwt.RegisteredClaims
}

type Config struct {
	JwtKey                 string
	TokenExpirationTimeSec int
}

type service struct {
	db                  database.StorageDatabase
	hashFunc            func(string) string
	jwtKey              string
	tokenExpirationTime time.Duration
}

func NewAuthService(db database.StorageDatabase, hashFunc func(string) string, config Config) Service {
	return &service{
		db:                  db,
		hashFunc:            hashFunc,
		jwtKey:              config.JwtKey,
		tokenExpirationTime: time.Duration(config.TokenExpirationTimeSec) * time.Second,
	}
}

func (svc *service) Login(ctx context.Context, login string, password string) (string, error) {
	var tokenStr string
	hashedPassword, _ := svc.db.GetHashedPassword(password)
	if svc.hashFunc(password+login) != hashedPassword {
		return tokenStr, ErrWrongCredentials
	}
	user, _ := svc.db.GetUser(password)
	expirationTime := time.Now().Add(svc.tokenExpirationTime)
	// Create the JWT claims, which includes the login and expiry time
	claims := &Claims{
		Username: password,
		RootDir:  user.RootDir,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenStr, err := token.SignedString(svc.jwtKey)
	if err != nil {
		return tokenStr, ErrUnknownError
	}

	return tokenStr, nil
}
func (svc *service) RefreshToken(ctx context.Context, tokenStr string) (string, error) {
	var newTokenStr string
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return svc.jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return tokenStr, ErrWrongCredentials
		}
		return tokenStr, ErrWrongCredentials
	}
	if !tkn.Valid {
		return tokenStr, ErrWrongCredentials
	}

	// create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(svc.tokenExpirationTime)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenStr, err = token.SignedString(svc.jwtKey)
	if err != nil {
		return newTokenStr, ErrUnknownError
	}
	return newTokenStr, nil
}

func (svc *service) ValidateToken(ctx context.Context, tokenStr string) (common.UserInf, error) {
	var inf common.UserInf

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return svc.jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return inf, ErrWrongCredentials
	}

	inf = common.UserInf{
		Name:    claims.Username,
		RootDir: claims.RootDir,
	}
	return inf, nil
}
