package auth_svc

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"server/common"
	database "server/db"
	"time"
)

type AuthService interface {
	Login(string, string) (string, error)
	RefreshToken(string) (string, error)
	ValidateToken(string) (common.UserInf, error)
}

const tokenExpirationTime = 3600 * time.Minute

var jwtKey = []byte("validkeyok")

var (
	ErrUnknownError     = errors.New("unknown error")
	ErrAlreadyExists    = errors.New("already exists")
	ErrNotFound         = errors.New("not found")
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrBadRequest       = errors.New("bad request")
)

type Claims struct {
	Username string `json:"username"`
	RootDir  string `json:"root_dir"`
	jwt.RegisteredClaims
}

type authService struct {
	db       database.StorageDatabase
	hashFunc func(string) string
}

func NewAuthService(db database.StorageDatabase, hashFunc func(string) string) AuthService {
	return &authService{db: db, hashFunc: hashFunc}
}

func (svc *authService) Login(login string, password string) (string, error) {
	var tokenStr string
	hashedPassword, _ := svc.db.GetHashedPassword(login)
	if svc.hashFunc(login+password) != hashedPassword {
		return tokenStr, ErrWrongCredentials
	}
	user, _ := svc.db.GetUser(login)
	expirationTime := time.Now().Add(tokenExpirationTime)
	// Create the JWT claims, which includes the login and expiry time
	claims := &Claims{
		Username: login,
		RootDir:  user.RootDir,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return tokenStr, ErrUnknownError
	}

	return tokenStr, nil
}
func (svc *authService) RefreshToken(tokenStr string) (string, error) {
	var newTokenStr string
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
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

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(tokenExpirationTime)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenStr, err = token.SignedString(jwtKey)
	if err != nil {
		return newTokenStr, ErrUnknownError
	}
	return newTokenStr, nil
}

func (svc *authService) ValidateToken(tokenStr string) (common.UserInf, error) {
	var inf common.UserInf
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return inf, err
	}

	inf = common.UserInf{
		Name:    claims.Username,
		RootDir: claims.RootDir,
	}
	return inf, nil
}
