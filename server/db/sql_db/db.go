package sql_db

import "database/sql"

// Database interface that provides Connect() and Close() methods and include Query interface
type Database interface {
	Connect() error
	Close() error
	Query
}

// Query interface provides methods for execution sql queries
type Query interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
