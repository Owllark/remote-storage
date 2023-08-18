package db

import "server/common"

type StorageDatabase interface {
	Connect(user, password, dbname, host string) error
	GetUsers() ([]common.UserInf, error)
	GetUser(name string) (common.UserInf, error)
	CreateUser(name, password, rootDir string) error
	GetHashedPassword(name string) (string, error)
}
