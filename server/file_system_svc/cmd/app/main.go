package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"server/file_system_svc"
	"strconv"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		Name     string `json:"name"`
		User     string `json:"user"`
	} `json:"database"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	ConsulServerAddress string `json:"consulServerAddress"`
	RootDirectory       string `json:"rootDirectory"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func main() {
	config := LoadConfiguration("file_system_svc\\config\\config.json")
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

	var s file_system_svc.FileSystemService
	{
		s = file_system_svc.NewFileSystemService(logger, file_system_svc.Config{
			RootDir:             config.RootDirectory,
			ConsulServerAddress: config.ConsulServerAddress,
		})
		s = file_system_svc.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = file_system_svc.MakeHttpHandler(s)
	}
	host := config.Host
	port := config.Port
	address := host + ":" + strconv.Itoa(port)
	http.ListenAndServe(address, h)
}
