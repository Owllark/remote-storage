package db

import "remote-storage/server/src/helper"

type StorageDatabase interface {
	Connect(user, password, dbname, host string) error
	GetClients(name string) ([]helper.UserInf, error)
	CreateClient(name, password string) error
}
