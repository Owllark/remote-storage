package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"remote-storage/server/src/db"
	"remote-storage/server/src/file_system"
	"remote-storage/server/src/helper"
	"strings"
)

var rootFs file_system.FileSystem

const pathSeparator = string(os.PathSeparator)

type uploadState struct {
	chunksNum      int
	tempDir        string
	fileName       string
	fileLocation   string
	chunksGotten   int
	receivedChunks []bool
}

type downloadState struct {
	chunks [][]byte
}

type User struct {
	upload   uploadState
	download downloadState
	inf      helper.UserInf
	fs       file_system.FileSystem
	//rootDir  string
}

var users []User

var database db.StorageDatabasePG

const (
	dbUser     = "postgres"
	dbHost     = "localhost"
	dbName     = "remote_storage"
	dbPassword = "password"
)

var jwtKey = []byte("secret")

const storagePath = "storage" + pathSeparator

func main() {

	err := database.Connect(dbUser, dbPassword, dbName, dbHost)
	if err != nil {
		log.Fatal(err)
	}

	ClientsInit()

	rootFs.SetRootDir("storage" + pathSeparator)
	r := mux.NewRouter()

	r.HandleFunc("/state", GetState)

	r.HandleFunc("/rename", Rename)
	r.HandleFunc("/move", Move)
	r.HandleFunc("/copy", Copy)
	r.HandleFunc("/delete", Delete)

	r.HandleFunc("/mkdir", MkDir)

	r.HandleFunc("/upload", StartUploading)
	r.HandleFunc("/upload/chunk", UploadChunk)
	r.HandleFunc("/upload/completed", UploadComplete)

	r.HandleFunc("/download", StartDownloading)
	r.HandleFunc("/download/chunk", DownloadChunk)

	r.HandleFunc("/authenticate", Authenticate)
	r.HandleFunc("/refresh", Refresh)
	r.HandleFunc("/logout", Logout)

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
				inf: helper.UserInf{Name: client.Name},
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
				clientInf := helper.UserInf{
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
					inf: helper.UserInf{Name: clientInf.Name},
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
