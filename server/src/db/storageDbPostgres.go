package db

import (
	_ "github.com/lib/pq"
	"remote-storage/server/src/db/sql_db"
	"remote-storage/server/src/helper"
)

// StorageDatabasePG provides higher level of database abstraction
// contains sql_db.PostgreSQLDatabase and provides methods for necessary queries
type StorageDatabasePG struct {
	db sql_db.PostgreSQLDatabase
}

// Connect takes arguments user, password, dbname, host for connecting to database
// load migration files and executes migrations
// returns error of database connection
func (s *StorageDatabasePG) Connect(user, password, dbname, host string) error {
	/*err := s.db.ConnectToHost(user, password, host)
	if err != nil {
		return err
	}*/
	/*err = MigrationsUp(s.db.Conn, "db/create_db_migrations")
	if err != nil {
		return err
	}*/
	err := s.db.Connect(user, password, dbname, host)
	if err != nil {
		return err
	}
	/*err = MigrationsUp(s.db.Conn, "db/migrations")*/
	return err
}

func (s *StorageDatabasePG) GetHashedPassword(name string) (string, error) {
	var res string
	var hashedPassword []uint8
	row := s.db.QueryRow(
		"SELECT hashed_password FROM clients WHERE name=$1",
		name,
	)

	err := row.Scan(&hashedPassword)
	res = string(hashedPassword)
	return res, err
}

func (s *StorageDatabasePG) GetUsers() ([]helper.UserInf, error) {
	var res []helper.UserInf
	var err error

	rows, err := s.db.Query(
		"SELECT name, root_directory FROM clients",
	)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var name []uint8
		var rootDir []uint8
		err := rows.Scan(&name, &rootDir)
		if err != nil {
			continue
		}
		inf := helper.UserInf{
			Name:    string(name),
			RootDir: string(rootDir),
		}
		res = append(res, inf)
	}
	err = rows.Err()

	return res, err
}

func (s *StorageDatabasePG) CreateUser(name, password, rootDir string) error {
	_, err := s.db.Exec(
		"INSERT INTO clients (name, hashed_password, root_directory) VALUES($1, $2, $3)",
		name,
		helper.Hash(name+password),
		rootDir,
	)
	return err
}
