package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"remote-storage/server/api"
	"remote-storage/server/db"
	"remote-storage/server/utils"
	"strings"
)

var rootFs file_system.FileSystem

const pathSeparator = string(os.PathSeparator)

var users []User

var database db.StorageDatabasePG

const (
	dbUser     = "postgres"
	dbHost     = "localhost"
	dbName     = "remote_storage"
	dbPassword = "password"
)

const storagePath = "storage" + pathSeparator

func main() {

	err := database.Connect(dbUser, dbPassword, dbName, dbHost)
	if err != nil {
		log.Fatal(err)
	}

	ClientsInit()

	rootFs.SetRootDir("storage" + pathSeparator)
	r := mux.NewRouter()

	r.HandleFunc("/state", api.GetState)

	r.HandleFunc("/rename", api.Rename)
	r.HandleFunc("/move", api.Move)
	r.HandleFunc("/copy", api.Copy)
	r.HandleFunc("/delete", api.Delete)

	r.HandleFunc("/mkdir", api.MkDir)

	r.HandleFunc("/upload", api.StartUploading)
	r.HandleFunc("/upload/chunk", api.UploadChunk)
	r.HandleFunc("/upload/completed", api.UploadComplete)

	r.HandleFunc("/download", api.StartDownloading)
	r.HandleFunc("/download/chunk", api.DownloadChunk)

	r.HandleFunc("/authenticate", api.Authenticate)
	r.HandleFunc("/refresh", api.Refresh)
	r.HandleFunc("/logout", api.Logout)

	go AcceptInput()
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

}

func ClientsInit() {
	clientList, err := database.GetUsers()
	if err != nil {
		log.Fatal("Failed to initialize users:", err)
	} else {
		for _, client := range clientList {
			var fs file_system.FileSystem
			fs.SetRootDir(client.RootDir)
			newClient := User{
				inf: utils.UserInf{Name: client.Name},
				fs:  fs,
			}
			users = append(users, newClient)
		}
	}
}

func AcceptInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		var req string
		req, _ = reader.ReadString('\n')
		req = req[0 : len(req)-2]
		arguments := strings.Split(req, " ")
		switch arguments[0] {
		case "newuser":
			{
				if len(arguments) < 3 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 3 {
					fmt.Println("Too much arguments")
					continue
				}
				name := arguments[1]
				password := arguments[2]
				clientInf := utils.UserInf{
					Name:    name,
					RootDir: storagePath + name + pathSeparator,
				}
				err := database.CreateUser(name, password, clientInf.RootDir)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("User created successfully")
				}
				rootFs.MkDir("", clientInf.Name)
				var fs file_system.FileSystem
				fs.SetRootDir(clientInf.RootDir)
				newClient := User{
					inf: utils.UserInf{Name: clientInf.Name},
					fs:  fs,
				}
				users = append(users, newClient)

			}
		case "getusers":
			{
				if len(arguments) < 1 {
					fmt.Println("Not enough arguments")
					continue
				} else if len(arguments) > 1 {
					fmt.Println("Too much arguments")
					continue
				}
				userNames, err := database.GetUsers()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Users:")
					for _, c := range userNames {
						fmt.Println(c.Name)
					}
				}

			}
		default:
			{

			}

		}
	}
}
