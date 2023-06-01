package sql_db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgreSQLDatabase struct {
	Conn *sql.DB
	PostgreSQLQuery
}

func (p *PostgreSQLDatabase) ConnectToHost(user, password, host string) error {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=5432 user=%s password=%s sslmode=disable", host, user, password))
	if err != nil {
		return err
	}
	p.Conn = db
	p.queryConn = db
	return nil
}

func (p *PostgreSQLDatabase) Connect(user, password, dbname, host string) error {
	//sql.Open("postgres", fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", "localhost", "postgres", "password", "remote_storage"))

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname))
	if err != nil {
		return err
	}
	p.Conn = db
	p.queryConn = db
	return nil
}

func (p *PostgreSQLDatabase) Close() error {
	return p.Conn.Close()
}

type PostgreSQLQuery struct {
	queryConn *sql.DB
}

func (p *PostgreSQLQuery) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.queryConn.Exec(query, args...)
}

func (p *PostgreSQLQuery) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.queryConn.Query(query, args...)
}

func (p *PostgreSQLQuery) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.queryConn.QueryRow(query, args...)
}
