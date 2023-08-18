package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	apiconsul "github.com/hashicorp/consul/api"
	"net/http"
	"os"
	"server/authsvc"
	"server/common"
	"server/db"
	"strconv"
	"time"
)

var database db.StorageDatabasePG

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		Name     string `json:"name"`
		User     string `json:"user"`
	} `json:"database"`
	Host                   string `json:"host"`
	Port                   int    `json:"port"`
	JwtKey                 string `json:"jwtKey"`
	TokenExpirationTimeSec int    `json:"tokenExpirationTimeSec"`
	ConsulServerAddress    string `json:"consulServerAddress"`
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

func registerInConsul(serviceID, serviceName, address string, port int, tags []string) error {
	// Create a Consul client
	consulConfig := apiconsul.DefaultConfig()
	consulConfig.Address = "http://localhost:8500" // Replace with your Consul server address
	client, err := apiconsul.NewClient(consulConfig)
	if err != nil {
		return err
	}

	// Register the service with Consul
	registration := &apiconsul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Tags:    tags,
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	config := LoadConfiguration("authsvc\\config\\config.json")

	logFile, err := os.OpenFile("authsvc\\log\\log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(logFile)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	err = database.Connect(config.Database.User, config.Database.Password, config.Database.Name, config.Database.Host)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
	}

	var s authsvc.Service
	{
		s = authsvc.NewAuthService(&database, common.SHA256hashing, authsvc.Config{
			JwtKey:                 config.JwtKey,
			TokenExpirationTimeSec: config.TokenExpirationTimeSec,
		})
		s = authsvc.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = authsvc.MakeHttpHandler(s)
	}
	serviceID := "auth_service-1"
	serviceName := "auth-service"
	address := config.Host
	port := config.Port
	tags := []string{"prod"}
	registerInConsul(serviceID, serviceName, address, port, tags)

	http.ListenAndServe(address+":"+strconv.Itoa(port), h)
	logger.Log("ServiceSessionEnd", time.Now())
}
