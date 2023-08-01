package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"server/auth_svc"
	"server/common"
	"server/db"
	"time"
)

var database db.StorageDatabasePG

const (
	dbUser     = "postgres"
	dbHost     = "localhost"
	dbName     = "remote_storage"
	dbPassword = "password"
)

func main() {
	logFile, err := os.OpenFile("file_system_svc\\log\\log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(logFile)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	err = database.Connect(dbUser, dbPassword, dbName, dbHost)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
	}

	var s auth_svc.AuthService
	{
		s = auth_svc.NewAuthService(&database, common.Hash256)
		s = auth_svc.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = auth_svc.MakeHttpHandler(s)
	}

	http.ListenAndServe("localhost:666", h)
	logger.Log("ServiceSessionEnd", time.Now())
}
